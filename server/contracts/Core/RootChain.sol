pragma solidity ^0.4.22;

import 'openzeppelin-solidity/contracts/token/ERC721/ERC721Receiver.sol';
import 'openzeppelin-solidity/contracts/math/SafeMath.sol';
import './Cards.sol';

// Lib deps
import '../Queue/PriorityQueue.sol';
import '../Libraries/ERC721PlasmaRLP.sol';
import '../Libraries/ECVerify.sol';
import '../Libraries/Merkle.sol';

contract RootChain is ERC721Receiver {
    /*
     * Events
     */
    event Deposit(uint256 slot, uint256 depositBlockNumber, uint256 denomination, address indexed from);
    event CanWithdraw(address  owner, uint  uid);
    event FinalizedExit(uint priority, address  owner, uint256  uid);

    using SafeMath for uint256;
    using ERC721PlasmaRLP for bytes;
    using ERC721PlasmaRLP for ERC721PlasmaRLP.txData;
    using ECVerify for bytes32;
    using Merkle for bytes32;

    /*
     * Storage
     */

    address public authority;

    // exits
    PriorityQueue exitsQueue;
    mapping(uint256 => Exit) public exits;
    struct Exit {
        address owner;
        uint256 slot;
        uint256 created_at;
    }

    // tracking of NFTs deposited in each slot
    uint public NUM_COINS;
    mapping (uint => NFT_UTXO) public coins; 
    struct NFT_UTXO {
        uint256 uid; // there are up to 2^256 cards, can probably make it less
        uint256 denomination; // an owner cannot own more than 256 of a card. Currently set to 1 always, subject to change once the token changes
        address owner; // who owns that nft
        bool canExit;
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
    mapping(address => uint256[]) public pendingWithdrawals; //  the pending cards to withdraw
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

        exitsQueue = new PriorityQueue();
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
    function deposit(address from, uint256 uid, uint256 denomination)
        private
    {
        // Update state "tree"
        coins[NUM_COINS] = NFT_UTXO({
                uid: uid, 
                denomination: denomination,
                owner: from, 
                canExit: false // cannot directly withdraw a coin, need to go through normal exit process
            });
        // TX hash is always the hash slot of the coin
        bytes32 txHash = keccak256(NUM_COINS); 
        uint256 depositBlockNumber = getDepositBlock();

        childChain[depositBlockNumber] = childBlock({
            root: txHash,
            created_at: block.timestamp
        });

        currentDepositBlock = currentDepositBlock.add(1);
        emit Deposit(NUM_COINS, depositBlockNumber, denomination, from); // create a utxo at `uid`

        NUM_COINS += 1;
    }

    /// txBlk is not needed for sure, maybe to prevent replay attacks, still unclear. May help for splitting/merging

    function startExit(
        bytes prevTxBytes, bytes exitingTxBytes, 
        bytes prevTxInclusionProof, bytes exitingTxInclusionProof, 
        bytes exitingTxSig ) 
        external
    {
        ERC721PlasmaRLP.txData memory exitingTxData = exitingTxBytes.createExitingTx();

        //  Need to check that the exiting transaction has a valid signatrue by its owner in order to prevent someone else exiting the owner's funds when they don't want it.
        bytes32 txHash = keccak256(exitingTxData.slot);
        require(txHash.ecverify(exitingTxSig, exitingTxData.owner), "Invalid sig");

        if (exitingTxData.prevBlock % childBlockInterval != 0 ) { 
            // If it's an exit of a deposit transaction then we need to check that:
            // the transaction hash was indeed the root of the claimed deposit block
            require(txHash == childChain[exitingTxData.prevBlock].root, 
                    "Deposit Tx not included in block");
        } else {
            // Otherwise we need to check for a proper inclusion of the transaction and the
            // referenced transaction in the blocks' merkle roots
            require(checkBlockInclusion(
                    prevTxBytes, exitingTxBytes,
                    prevTxInclusionProof, exitingTxInclusionProof
                )
            );
        }
        uint priority = exitingTxData.prevBlock * 10000000  + exitingTxData.slot * 10000;

        exitsQueue.insert(priority);
        exits[priority] = Exit({
            owner: exitingTxData.owner, 
            slot: exitingTxData.slot, 
            created_at: block.timestamp
        });
    }

    function checkBlockInclusion(
            bytes prevTxBytes, bytes exitingTxBytes,
            bytes prevTxInclusionProof, bytes exitingTxInclusionProof) 
            returns (bool)
    {
        ERC721PlasmaRLP.txData memory prevTxData = prevTxBytes.createExitingTx();
        bytes32 prevMerkleHash = keccak256(prevTxBytes);
        bytes32 prevRoot = childChain[prevTxData.prevBlock].root;


        ERC721PlasmaRLP.txData memory exitingTxData = exitingTxBytes.createExitingTx();
        bytes32 merkleHash = keccak256(exitingTxBytes);
        bytes32 root = childChain[exitingTxData.prevBlock].root;

        require(
            prevMerkleHash.checkMembership(
                prevTxData.slot,
                prevRoot, 
                prevTxInclusionProof
            ),
            "Previous tx not included in claimed block"
        );

        require(
            merkleHash.checkMembership(
                exitingTxData.slot, 
                root, 
                exitingTxInclusionProof
            ),
            "Exiting tx not included in claimed block"
        );

        return true;
    }


    function finalizeExits() public {
        require(exitsQueue.currentSize() > 0, "exit queue empty");

        uint256 priority = exitsQueue.getMin();
        Exit memory currentExit = exits[priority];

        // finalize exits that are older than `1 week`.
        while (exitsQueue.currentSize() > 0 && 
               (block.timestamp - currentExit.created_at) > 1 weeks) {
            // this can occur if challengeExit is sucessful on an exit
            if (currentExit.owner == address(0)) { // handles exits
                exitsQueue.delMin();
                if (exitsQueue.currentSize() == 0) return; // no revert because we wwant to keep the delete

                // move onto the next oldest exit
                priority = exitsQueue.getMin();
                currentExit = exits[priority];
                continue; // Prevent incorrect processing of deleted exits.
            }

            // Change owner of coin at exit.slot
            coins[currentExit.slot].owner = currentExit.owner;
            coins[currentExit.slot].canExit = true;

            emit FinalizedExit(priority, currentExit.owner, currentExit.slot);
            emit CanWithdraw(currentExit.owner, currentExit.slot);

            // delete the finalized exit
            exitsQueue.delMin();
            delete exits[priority];

            // move onto the next oldest exit
            if (exitsQueue.currentSize() == 0) {
                return;
            }
            priority = exitsQueue.getMin();
            currentExit = exits[priority];
        }
    }

    // Withdraw a UTXO that has been exited
    function withdraw(uint slot) external {
        require(coins[slot].owner == msg.sender, "You do not own that UTXO");
        require(coins[slot].canExit, "You cannot exit that coin!");
        cryptoCards.safeTransferFrom(address(this), msg.sender, coins[slot].uid);
    }

    function getDepositBlock() public view returns (uint256) {
        return currentChildBlock.sub(childBlockInterval).add(currentDepositBlock);
    }

    /// receiver for erc721 to trigger a deposit
    function onERC721Received(address _from, uint256 _uid, bytes) 
        public 
        returns(bytes4) 
    {
        require(msg.sender == address(cryptoCards)); // can only be called by the associated cryptocards contract. 
        deposit(_from, _uid, 1); //, _data);
        return ERC721_RECEIVED;
    }



}

