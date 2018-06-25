pragma solidity ^0.4.24;

// Zeppelin Imports
import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Receiver.sol";
import "openzeppelin-solidity/contracts/math/SafeMath.sol";

// Lib deps
import "../Libraries/Transaction/Transaction.sol";
import "../Libraries/ECVerify.sol";

import "./SparseMerkleTree.sol";
import "./ValidatorManagerContract.sol";


contract RootChain is ERC721Receiver {
    event Deposit(uint64 indexed slot, uint256 blockNumber, uint64 denomination, address indexed from);
    event SubmittedBlock(uint256 blockNumber, bytes32 root, uint256 timestamp);

    event StartedExit(uint64 indexed slot, address indexed owner);
    event ChallengedExit(uint64 indexed slot);
    event RespondedExitChallenge(uint64 indexed slot);
    event FinalizedExit(uint64  indexed slot, address owner);

    event FreedBond(address indexed from, uint256 amount);
    event SlashedBond(address indexed from, address indexed to, uint256 amount);
    event WithdrewBonds(address indexed from, uint256 amount);

    using SafeMath for uint256;
    using Transaction for bytes;
    using ECVerify for bytes32;

    uint256 constant BOND_AMOUNT = 0.1 ether;

    address public authority;

    /*
     * Modifiers
     */
    modifier isValidator() {
        require(vmc.checkValidator(msg.sender));
        _;
    }

    modifier isTokenApproved(address _address) {
        require(vmc.allowedTokens(_address));
        _;
    }

    modifier isBonded() {
        require(msg.value == BOND_AMOUNT);

        // Save challenger's bond
        balances[msg.sender].bonded = balances[msg.sender].bonded.add(msg.value);
        _;
    }

    modifier isState(uint64 slot, State state) {
        require(coins[slot].state == state, "Wrong state");
        _;
    }

    modifier cleanupExit(uint64 slot) {
        _;
        delete coins[slot].exit;
        delete exitSlots[getIndex(slot)];
    }

    struct Balance {
        uint256 bonded;
        uint256 withdrawable;
    }
    mapping (address => Balance) public balances;

    // exits
    uint64[] public exitSlots;
    // Each exit can only be challenged by a single challenger at a time
    mapping (uint64 => address) challengers;
    struct Exit {
        address prevOwner; // previous owner of coin
        address owner;
        uint256 createdAt;
        uint256 bond;
        uint256 prevBlock;
        uint256 exitBlock;
    }
    enum State {
        DEPOSITED,
        EXITING,
        CHALLENGED,
        RESPONDED,
        EXITED
    }

    // tracking of NFTs deposited in each slot
    uint64 public numCoins = 0;
    mapping (uint64 => Coin) public coins;
    struct Coin {
        uint64 uid; // there are up to 2^64 cards, can probably make it less
        uint32 denomination; // Currently set to 1 always, subject to change once the token changes
        uint256 depositBlock;
        address owner; // who owns that nft
        address contractAddress; // which contract does the coin belong to
        State state;
        Exit exit;
    }

    // child chain
    uint256 public childBlockInterval = 1000;
    uint256 public currentBlock = 0;
    struct ChildBlock {
        bytes32 root;
        uint256 createdAt;
    }

    uint256 public depositCount;
    mapping(uint256 => ChildBlock) public childChain;
    ValidatorManagerContract vmc;
    SparseMerkleTree smt;

    constructor (ValidatorManagerContract _vmc) public {
        vmc = _vmc;
        smt = new SparseMerkleTree();
    }

    function submitBlock(bytes32 root)
        public
        isValidator
    {
        // rounding to next whole `childBlockInterval`
        currentBlock = currentBlock.add(childBlockInterval)
            .div(childBlockInterval)
            .mul(childBlockInterval);

        childChain[currentBlock] = ChildBlock({
            root: root,
            createdAt: block.timestamp
        });

        emit SubmittedBlock(currentBlock, root, block.timestamp);
    }

    /// @dev Allows anyone to deposit funds into the Plasma chain, called when contract receives ERC721
    function deposit(address from, uint64 uid, uint32 denomination)
        private
    {
        currentBlock = currentBlock.add(1);

        // Update state. Leave `exit` empty
        Coin memory coin;
        coin.uid = uid;
        coin.contractAddress = msg.sender;
        coin.denomination = denomination;
        coin.depositBlock = currentBlock;
        coin.owner = from;
        coin.state = State.DEPOSITED;
        uint64 slot = uint64(bytes8(keccak256(abi.encodePacked(numCoins, msg.sender, from))));
        coins[slot] = coin;

        childChain[currentBlock] = ChildBlock({
            // save signed transaction hash as root
            // hash for deposit transactions is the hash of its slot
            root: keccak256(abi.encodePacked(slot)),
            createdAt: block.timestamp
        });

        // create a utxo at `slot`
        emit Deposit(
            slot,
            currentBlock,
            denomination,
            from);

        numCoins += 1;
    }

    /******************** EXIT RELATED ********************/

    function startExit(
        uint64 slot, bytes prevTxBytes,
        bytes exitingTxBytes, bytes prevTxInclusionProof,
        bytes exitingTxInclusionProof, bytes sig,
        uint256 prevTxIncBlock, uint256 exitingTxIncBlock)
        external
        payable isBonded
        isState(slot, State.DEPOSITED)
    {
        // If we're exiting a deposit UTXO, we do a different inclusion check
        if (exitingTxIncBlock % childBlockInterval != 0) {
            checkDepositBlockInclusion(exitingTxBytes, sig, exitingTxIncBlock, true);
        } else {
            checkBlockInclusion(
                prevTxBytes, exitingTxBytes,
                prevTxInclusionProof, exitingTxInclusionProof,
                sig,
                prevTxIncBlock, exitingTxIncBlock,
                true
            );
        }
        pushExit(slot, prevTxBytes, prevTxIncBlock, exitingTxIncBlock);
    }

    function finalizeExit(uint64 slot) public {
        Coin storage coin = coins[slot];

        // If a coin is not under exit/challenge, then ignore it
        if (coin.state == State.DEPOSITED || coin.state == State.EXITED)
            return;

        // If an exit is not matured, ignore it
        if ((block.timestamp - coin.exit.createdAt) <= 7 days)
            return;

        // If a coin has been challenged AND not responded, slash it
        if (coin.state == State.CHALLENGED) {
            // Update coin state & penalize exitor
            coin.state = State.DEPOSITED;
            slashBond(coin.exit.owner, challengers[slot]);
        // otherwise, the exit has not been challenged, or it has been challenged and responded
        } else {
            // If the exit was actually challenged and responded, penalize the challenger
            if (coin.state == State.RESPONDED) {
                slashBond(challengers[slot], coin.exit.owner);
            }

            // Update coin's owner
            coin.owner = coin.exit.owner;
            coin.state = State.EXITED;

            // Allow the exitor to withdraw their bond
            freeBond(coin.owner);

            emit FinalizedExit(slot, coin.owner);
        }
        delete coins[slot].exit;
        delete exitSlots[getIndex(slot)];
    }

    function finalizeExits() external {
        uint256 exitSlotsLength = exitSlots.length;
        for (uint256 i = 0; i < exitSlotsLength; i++) {
            finalizeExit(exitSlots[i]);
        }
    }

    // Withdraw a UTXO that has been exited
    function withdraw(uint64 slot) external isState(slot, State.EXITED) {
        require(coins[slot].owner == msg.sender, "You do not own that UTXO");
        ERC721(coins[slot].contractAddress).safeTransferFrom(address(this), msg.sender, uint256(coins[slot].uid));
    }

    /******************** CHALLENGES ********************/

    // Submit proof of a transaction before prevTx
    // Exitor has to call respondChallengeBefore and submit a transaction before prevTx or prevTx itself.
    function challengeBefore(
        uint64 slot,
        bytes prevTxBytes, bytes exitingTxBytes,
        bytes prevTxInclusionProof, bytes exitingTxInclusionProof,
        bytes sig,
        uint256 prevTxIncBlock, uint256 exitingTxIncBlock)
        external
        payable isBonded
        isState(slot, State.EXITING)
    {
        require(prevTxIncBlock < exitingTxIncBlock);
        // If we're exiting a deposit UTXO directly, we do a different inclusion check
        if (exitingTxIncBlock % childBlockInterval != 0) {
            checkDepositBlockInclusion(
                exitingTxBytes,
                sig,
                exitingTxIncBlock,
                false);
        } else {
            checkBlockInclusion(
                prevTxBytes, exitingTxBytes, prevTxInclusionProof,
                exitingTxInclusionProof, sig,
                prevTxIncBlock, exitingTxIncBlock, false);
        }
        setChallenged(slot);
    }

    // If `challengeBefore` was successfully challenged, then set state to RESPONDED and allow the coin to be exited. No need to actually attach a bond when responding to a challenge
    function respondChallengeBefore(
        uint64 slot,
        uint256 challengingBlockNumber,
        bytes challengingTransaction,
        bytes proof)
        external
        isState(slot, State.CHALLENGED)
    {
        checkTxIncluded(challengingTransaction, challengingBlockNumber, proof);

        // Mark exit as responded, which will allow it to be finalized once it has matured
        coins[slot].state = State.RESPONDED;
        emit RespondedExitChallenge(slot);
    }

    function challengeBetween(
        uint64 slot,
        uint256 challengingBlockNumber,
        bytes challengingTransaction,
        bytes proof,
        bytes signature)
        external isState(slot, State.EXITING) cleanupExit(slot)
    {
        // Must challenge with a tx in between
        require(
            coins[slot].exit.exitBlock > challengingBlockNumber && coins[slot].exit.prevBlock < challengingBlockNumber,
            "Challenging transaction must have happened AFTER the attested exit's timestamp");

        // Check that the challenging transaction has been signed
        // by the attested previous owner of the coin in the exit
        bytes32 txHash = keccak256(challengingTransaction);
        require(txHash.ecverify(signature, coins[slot].exit.prevOwner), "Invalid sig");

        checkTxIncluded(challengingTransaction, challengingBlockNumber, proof);
        // Apply penalties and change state
        slashBond(coins[slot].exit.owner, msg.sender);
        coins[slot].state = State.DEPOSITED;
    }

    function challengeAfter(
        uint64 slot,
        uint256 challengingBlockNumber,
        bytes challengingTransaction,
        bytes proof,
        bytes signature)
        external
        isState(slot, State.EXITING)
        cleanupExit(slot)
    {
        checkDirectSpend(slot, challengingTransaction, signature);
        checkTxIncluded(challengingTransaction, challengingBlockNumber, proof);
        // Apply penalties and delete the exit
        slashBond(coins[slot].exit.owner, msg.sender);
        // Reset coin state
        coins[slot].state = State.DEPOSITED;
    }

    /******************** BOND RELATED ********************/

    function freeBond(address from) private {
        balances[from].bonded = balances[from].bonded.sub(BOND_AMOUNT);
        balances[from].withdrawable = balances[from].withdrawable.add(BOND_AMOUNT);
        emit FreedBond(from, BOND_AMOUNT);
    }

    function withdrawBonds() external {
        // Can only withdraw bond if the msg.sender
        uint256 amount = balances[msg.sender].withdrawable;
        balances[msg.sender].withdrawable = 0; // no reentrancy!

        msg.sender.transfer(amount);
        emit WithdrewBonds(msg.sender, amount);
    }

    function slashBond(address from, address to) private {
        balances[from].bonded = balances[from].bonded.sub(BOND_AMOUNT);
        balances[to].withdrawable = balances[to].withdrawable.add(BOND_AMOUNT);
        emit SlashedBond(from, to, BOND_AMOUNT);
    }

    function pushExit(
        uint64 slot,
        bytes txBytes,
        uint256 prevBlock,
        uint256 exitingBlock) private
    {
        Transaction.TX memory txData = txBytes.getTx();

        // Push exit to list
        exitSlots.push(slot);

        // Create exit
        Coin storage c = coins[slot];
        c.exit = Exit({
            prevOwner: txData.owner,
            owner: msg.sender,
            createdAt: block.timestamp,
            bond: msg.value,
            prevBlock: prevBlock,
            exitBlock: exitingBlock
        });

        // Update coin state
        c.state = State.EXITING;
        emit StartedExit(slot, msg.sender);
    }

    function setChallenged(uint64 slot) private {
        // Do not delete exit yet. Set its state as challenged and wait for the exitor's response
        coins[slot].state = State.CHALLENGED;
        // Save the challenger's address, for applying penalties
        challengers[slot] = msg.sender;
        emit ChallengedExit(slot);
    }

    /******************** PROOF CHECKING ********************/

    function checkDirectSpend(uint64 slot, bytes txBytes, bytes signature) private view {
        Transaction.TX memory txData = txBytes.getTx();
        bytes32 txHash = keccak256(txBytes);
        require(txHash.ecverify(signature, coins[slot].exit.owner), "Invalid sig");
        require(txData.prevBlock == coins[slot].exit.exitBlock, "Not a direct spend");
    }

    function checkDepositBlockInclusion(
        bytes txBytes,
        bytes signature,
        uint256 txIncBlock,
        bool checkSender
    )
         private
         view
         returns (bool)
    {
        Transaction.TX memory txData = txBytes.getTx();
        if (checkSender)
            require(txData.owner == msg.sender, "Invalid sender");

        bytes32 txHash = keccak256(abi.encodePacked(txData.slot));
        require(txHash.ecverify(signature, txData.owner), "Invalid sig");
        require(
            txHash == childChain[txIncBlock].root,
            "Deposit Tx not included in block"
        );

        return true;
    }

    function checkBlockInclusion(
        bytes prevTxBytes, bytes exitingTxBytes,
        bytes prevTxInclusionProof, bytes exitingTxInclusionProof,
        bytes sig,
        uint256 prevTxIncBlock, uint256 exitingTxIncBlock, bool checkSender)
        private
        view
        returns (bool)
    {
        Transaction.TX memory exitingTxData = exitingTxBytes.getTx();
        Transaction.TX memory prevTxData = prevTxBytes.getTx();

        if (checkSender)
            require(exitingTxData.owner == msg.sender, "Invalid sender");
        require(keccak256(exitingTxBytes).ecverify(sig, prevTxData.owner), "Invalid sig");
        require(exitingTxData.slot == prevTxData.slot);
        require(prevTxIncBlock < exitingTxIncBlock);

        checkTxIncluded(prevTxBytes, prevTxIncBlock, prevTxInclusionProof);
        checkTxIncluded(exitingTxBytes, exitingTxIncBlock, exitingTxInclusionProof);

        return true;
    }

    function checkTxIncluded(bytes txBytes, uint256 blockNumber, bytes proof) private view {
        Transaction.TX memory txData = txBytes.getTx();
        bytes32 txHash;
        bytes32 root = childChain[blockNumber].root;

        // If deposit block, just check the matching hash to the block's root. Simpler verification due to deposit mechanism
        if (blockNumber % childBlockInterval != 0) {
            txHash = keccak256(abi.encodePacked(txData.slot));
            require(txHash == root);
        } else {
            txHash = keccak256(txBytes);
            require(
                checkMembership(
                    txHash,
                    root,
                    txData.slot,
                    proof
                ),
                "Tx not included in claimed block"
            );
        }
    }

    /******************** ERC721 ********************/

    function onERC721Received(address _from, uint256 _uid, bytes)
        public
        isTokenApproved(msg.sender)
        returns(bytes4)
    {
        deposit(_from, uint64(_uid), uint32(1));
        return ERC721_RECEIVED;
    }

    /******************** HELPERS ********************/

    function getIndex(uint64 slot) private view returns (uint256) {
        uint256 len = exitSlots.length;
        for (uint256 i = 0; i < len; i++) {
            if (exitSlots[i] == slot)
                return i;
        }
        return 0;
    }

    function checkMembership(
        bytes32 txHash,
        bytes32 root,
        uint64 slot,
        bytes proof) public returns (bool)
    {
        return smt.checkMembership(
            txHash,
            root,
            slot,
            proof);
    }

    function getPlasmaCoin(uint64 slot) external view returns(uint64, uint256, uint32, address, State) {
        Coin memory c = coins[slot];
        return (c.uid, c.depositBlock, c.denomination, c.owner, c.state);
    }

    function getExit(uint64 slot) external view returns(address, uint256, uint256, State) {
        Exit memory e = coins[slot].exit;
        return (e.owner, e.prevBlock, e.exitBlock, coins[slot].state);
    }

    function getBlockRoot(uint256 blockNumber) public view returns (bytes32 root) {
        root = childChain[blockNumber].root;
    }
}
