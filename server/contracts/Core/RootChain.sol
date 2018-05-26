pragma solidity ^0.4.24;

// Linked contract for withdrawals, import only safeTransferFrom interface for gas efficiency in the future
import './Cards.sol';

// Zeppelin Imports
import 'openzeppelin-solidity/contracts/token/ERC721/ERC721Receiver.sol';
import 'openzeppelin-solidity/contracts/math/SafeMath.sol';

// Lib deps
import '../Libraries/Transaction/Transaction.sol';
import '../Libraries/ByteUtils.sol';
import '../Libraries/ECVerify.sol';

// Sparse Merkle Tree functionalities
import './SparseMerkleTree.sol';

contract RootChain is ERC721Receiver, SparseMerkleTree {
    /*
     * Events
     */
    event Deposit(uint64 indexed slot, uint256 depositBlockNumber, uint64 denomination, address indexed from);
    event StartedExit(uint64 indexed slot, address indexed owner, uint created_at);
    event ChallengedExit(uint64 indexed slot);
    event FinalizedExit(uint64  indexed slot, address owner);
    event SlashedBond(uint64 indexed slot, address indexed from, address indexed to, uint amount);

    using SafeMath for uint256; // if few operations consider removing and doing asserts inline for less gas costs
    using Transaction for bytes;
    using ECVerify for bytes32;

    /*
     * Storage
     */

    address public authority;

    uint constant BOND_AMOUNT = 0.01 ether;
    modifier isBonded() { 
        require ( msg.value == BOND_AMOUNT ); 
        _;
    }

    // bond[bob][coinA] = X ether
    struct Balance {
        uint bonded; // staked as a bond
        uint withdrawable;
    }
    mapping (address => Balance ) balances;

    // exits
    uint64[] public exitSlots;

    mapping (uint64 => address) challengers;
    struct Exit {
        address owner;
        uint256 created_at;
        uint256 bond;
    }

    enum State {
        DEPOSITED,
        EXITING,
        CHALLENGED,
        RESPONDED,
        EXITED
    }

    // tracking of NFTs deposited in each slot
    uint64 public NUM_COINS;
    mapping (uint64 => NFT_UTXO) public coins; 
    struct NFT_UTXO {
        uint64 uid; // there are up to 2^256 cards, can probably make it less
        uint32 denomination; // an owner cannot own more than 256 of a card. Currently set to 1 always, subject to change once the token changes
        address owner; // who owns that nft
        State state;
        Exit exit;
    }

    // child chain
    uint256 public childBlockInterval;
    uint256 public currentChildBlock;
    uint256 public currentDepositBlock;
    uint256 public lastParentBlock;
    struct childBlock {
        bytes32 root;
        uint256 created_at;
    }

    uint public depositCount;
    mapping(uint => childBlock) public childChain;
    CryptoCards cryptoCards;

    /*
     * Modifiers
     */
    modifier isAuthority() {
        require(msg.sender == authority);
        _;
    }

    constructor () public{
        authority = msg.sender;

		childBlockInterval = 1000;
        currentChildBlock = childBlockInterval;
        currentDepositBlock = 1;
        lastParentBlock = block.number; // to ensure no chain reorgs

    }

    function setCryptoCards(CryptoCards _cryptoCards) isAuthority public {
        cryptoCards = _cryptoCards;
    }

    /// @param root 32 byte merkleRoot of ChildChain block
    /// @notice childChain blocks can only be submitted at most every 6 root chain blocks
    function submitBlock(bytes32 root)
        public
        isAuthority
    {
        // ensure finality on previous blocks before submitting another
        // require(block.number >= lastParentBlock.add(6)); // commented out while prototyping

        childChain[currentChildBlock] = childBlock({
            root: root,
            created_at: block.timestamp
        });

        currentChildBlock = currentChildBlock.add(childBlockInterval);
        currentDepositBlock = 1;
        lastParentBlock = block.number;
	}


    /// @dev Allows anyone to deposit funds into the Plasma chain, called when contract receives ERC721
    function deposit(address from, uint64 uid, uint32 denomination, bytes txBytes)
        private
    {
        Transaction.TX memory txData = txBytes.getTx();
        // Verify that the transaction data sent matches the coin data from ERC721
        require(txData.slot == NUM_COINS);
        require(txData.denomination == denomination);
        require(txData.owner == from);
        require(txData.prevBlock == 0);

        // Update state. Leave `exit` empty
        NFT_UTXO memory coin;
        coin.uid = uid;
        coin.denomination = denomination;
        coin.owner = from;
        coin.state = State.DEPOSITED;
        coins[NUM_COINS] = coin;

        bytes32 txHash = keccak256(txBytes);
        uint256 depositBlockNumber = getDepositBlock();

        childChain[depositBlockNumber] = childBlock({
            root: txHash, // save signed transaction hash as root
            created_at: block.timestamp
        });

        currentDepositBlock = currentDepositBlock.add(1);
        emit Deposit(NUM_COINS, depositBlockNumber, denomination, from); // create a utxo at slot `NUM_COINS`

        NUM_COINS += 1;
    }

    function startExit(
        uint64 slot,
        bytes prevTxBytes, bytes exitingTxBytes, 
        bytes prevTxInclusionProof, bytes exitingTxInclusionProof, 
        bytes sigs,
        uint prevTxIncBlock, uint exitingTxIncBlock) 
        payable isBonded
        external
    {
        // Different inclusion check depending on if we're exiting a deposit transaction or not
        if (exitingTxIncBlock % childBlockInterval != 0 ) { 
           require(
                checkDepositBlockInclusion(
                    exitingTxBytes, 
                    sigs, // for deposit blocks this is just a single sig
                    exitingTxIncBlock
                ),
                "Not included in deposit block"
            );
        } else {
            require(
                checkBlockInclusion(
                    prevTxBytes, exitingTxBytes,
                    prevTxInclusionProof, exitingTxInclusionProof,
                    sigs,
                    prevTxIncBlock, exitingTxIncBlock
                ), 
                "Not included in blocks"
            );
        }

        exitSlots.push(slot);

        // Update Coin's storage
        NFT_UTXO storage c = coins[slot];
        c.exit = Exit({
            owner: msg.sender, 
            created_at: block.timestamp,
            bond: msg.value
        });
        c.state = State.EXITING;

        emit StartedExit(slot, msg.sender, block.timestamp);
    }

    function getSig(bytes sigs, uint i) public pure returns(bytes) {
        return ByteUtils.slice(sigs, 66 * i,  66);
    }

    function checkDepositBlockInclusion(
        bytes txBytes,
        bytes signature,
        uint txIncBlock
    )
         private 
         view 
         returns (bool) 
    {
        Transaction.TX memory txData = txBytes.getTx();
        bytes32 txHash = keccak256(txBytes); 

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
            bytes sigs,
            uint prevTxIncBlock, uint exitingTxIncBlock) 
            private
            view
            returns (bool)
    {
        Transaction.TX memory prevTxData = prevTxBytes.getTx();
        Transaction.TX memory exitingTxData = exitingTxBytes.getTx();

        bytes32 txHash = keccak256(exitingTxBytes);
        bytes32 root = childChain[exitingTxIncBlock].root;

        require(txHash.ecverify(getSig(sigs, 1), prevTxData.owner), "Invalid sig");
        require(exitingTxData.owner == msg.sender, "Invalid sender");
        
        require(
            checkMembership(
                txHash,
                root, 
                exitingTxData.slot, 
                exitingTxInclusionProof
            ),
            "Exiting tx not included in claimed block"
        );

        bytes32 prevTxHash = keccak256(prevTxBytes);
        bytes32 prevRoot = childChain[prevTxIncBlock].root;

        if (prevTxIncBlock % childBlockInterval != 0 ) { 
            require(prevTxHash == prevRoot); // like in deposit block
        } else {
            require(
                checkMembership(
                    prevTxHash,
                    prevRoot, 
                    prevTxData.slot,
                    prevTxInclusionProof
                ),
                "Previous tx not included in claimed block"
            );
        }

        return true;
    }

    function finalizeExit(uint64 slot) private returns (bool exited) {
        NFT_UTXO storage coin = coins[slot];
         if ((block.timestamp - coin.exit.created_at) > 7 days) {
             if (coin.state == State.CHALLENGED) {
                // If a coin's state is CHALLENGED, it needs to be responded. If not responded, the exit was invalid and thus the exitor's deposit must be slashed. 
             } 
             else { // If responded or exiting, proceed with the exit
                coin.owner = coin.exit.owner;
                coin.state = State.EXITED;
             }
            emit FinalizedExit(slot, coin.owner);
            return true;
         } else {
             return false;
         }
    }

    function finalizeExits() external {
        uint exitSlotsLength = exitSlots.length;
        uint64 slot;
        for (uint i = 0; i < exitSlotsLength; i++) { 
            slot = exitSlots[i];
            NFT_UTXO storage coin = coins[slot];
            if (coin.state == State.DEPOSITED || coin.state == State.EXITED) continue; // If a coin is not under exit/challenge, then ignore it
            if (finalizeExit(slot)) {
                delete coin.exit;
                delete exitSlots[i];
            }
        }
    }


    // CHALLENGES Require bonds in order to avoid griefing! 

    // Submit proof of a transaction before prevTx 
    // Exitor has to call respondChallengeBefore and submit a transaction before prevTx or prevTx itself.
    function challengeBefore(uint64 slot, bytes challengingTransaction, bytes proof) external {
        require(coins[slot].state == State.EXITING, "Coin not being exited");
        challengingTransaction;
        proof;

        // Do not delete exit yet. Set its state as challenged and wait for the exitor's response
        challengers[slot] = msg.sender; 
        coins[slot].state = State.CHALLENGED;
        emit ChallengedExit(slot);
    }

    // If `challengeBefore` was successfully challenged, then set state to RESPONDED and allow the coin to be exited. 
    function respondChallengeBefore(uint64 slot, bytes challengingTransaction, bytes proof) external {
        require(coins[slot].state == State.CHALLENGED, "Coin not under challenge");
        challengingTransaction;
        proof;
        coins[slot].state = State.RESPONDED;
    }


    function challengeBetween(uint64 slot, bytes challengingTransaction, bytes proof) external {
        require(coins[slot].state == State.EXITING, "Coin not being exited");

        // Validate proofs: TODO
        challengingTransaction;
        proof;

        // Apply penalties and delete the exit
        slashBond(slot, coins[slot].exit.owner, msg.sender);
        delete coins[slot].exit;    
        delete exitSlots[slot];

        // Reset coin state
        coins[slot].state = State.DEPOSITED;
    }

    function challengeAfter(uint64 slot, bytes challengingTransaction, bytes proof) external {
        require(coins[slot].state == State.EXITING, "Coin not being exited");

        // Validate proofs: TODO
        challengingTransaction;
        proof;

        // Apply penalties and delete the exit
        slashBond(slot, coins[slot].exit.owner, msg.sender);
        delete coins[slot].exit;    
        delete exitSlots[slot];

        // Reset coin state
        coins[slot].state = State.DEPOSITED;
    }

    function slashBond(uint64 slot, address from, address to) {
        balances[from].bonded = balances[from].bonded.sub(BOND_AMOUNT);
        balances[to].withdrawable = balances[to].withdrawable.add(BOND_AMOUNT);
        emit SlashedBond(slot, from, to, BOND_AMOUNT);
    
    }

    function withdrawBalance() external {
        // Can only withdraw bond if the msg.sender 
        uint amount = balances[msg.sender].withdrawable;
        balances[msg.sender].withdrawable = 0; // no reentrancy:
        msg.sender.transfer(amount);
    }


    // Withdraw a UTXO that has been exited
    function withdraw(uint64 slot) external {
        require(coins[slot].owner == msg.sender, "You do not own that UTXO");
        require(coins[slot].state == State.EXITED, "You cannot exit that coin!");
        cryptoCards.safeTransferFrom(address(this), msg.sender, uint256(coins[slot].uid));
    }

    function getDepositBlock() public view returns (uint256) {
        return currentChildBlock.sub(childBlockInterval).add(currentDepositBlock);
    }

    /// receiver for erc721 to trigger a deposit
    function onERC721Received(address _from, uint256 _uid, bytes _data) 
        public 
        returns(bytes4) 
    {
        require(msg.sender == address(cryptoCards)); // can only be called by the associated cryptocards contract. 
        deposit(_from, uint64(_uid), uint32(1), _data);
        return ERC721_RECEIVED;
    }
}

