pragma solidity ^0.4.22;

import 'zeppelin-solidity/contracts/token/ERC721/ERC721Receiver.sol';
import 'zeppelin-solidity/contracts/math/SafeMath.sol';
import './Queue/PriorityQueue.sol';
import './Cards.sol';
import './ERC721PlasmaRLP.sol';
import './ECVerify.sol';

contract RootChain is ERC721Receiver {
    /*
     * Events
     */
    event Deposit(address indexed depositor, uint256 tokenId, bytes data);
    event CanWithdraw(address  owner, uint  tokenId);
    event FinalizedExit(uint priority, address  owner, uint256  tokenId);

    using SafeMath for uint256;
    using ERC721PlasmaRLP for bytes;
    using ERC721PlasmaRLP for ERC721PlasmaRLP.txData;
    using ECVerify for bytes32;

    /*
     * Storage
     */

    address public authority;

    // exits
    PriorityQueue exitsQueue;
    mapping(uint256 => Exit) public exits;
    struct Exit {
        address owner;
        uint tokenId;
        uint256 created_at;
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
        require(block.number >= lastParentBlock.add(6));

        childChain[currentChildBlock] = childBlock({
            root: root,
            created_at: block.timestamp
        });

        currentChildBlock = currentChildBlock.add(childBlockInterval);
        currentDepositBlock = 1;
        lastParentBlock = block.number;
	}


    /// @dev Allows anyone to deposit funds into the Plasma chain, called when contract receives ERC721
    function deposit(address from, uint uid, bytes _data)
        private
    {
        // TODO: Serialize, do whatever with _data for UTXO/ChainID transfer
        // ERC721PlasmaRLP.txData memory depositTx = _data.createExitingTx(3);
        // uint prevBlock = depositTx.prevBlock;
        // uint uid = depositTx.uid;
        // address owner = depositTx.owner;

        bytes32 root = keccak256(_data);
        uint256 position = getDepositBlock();

        childChain[position] = childBlock({
            root: root,
            created_at: block.timestamp
        });

        currentDepositBlock = currentDepositBlock.add(1);
        emit Deposit(from, uid, _data); // create a utxo at `uid`
    }

    // Function still WIP
    /// Concept: Pass in the transaction that is being exited along with a reference to a previous valid transaction
    /// IF the previous transaction has valid merkle proof and was included in the specified block , then check that the signature on the previous transaction's new_owner is valid for the ecrecover for the current transaction. if valid, add to texits 
    // https://github.com/FourthState/plasma-mvp-rootchain/blob/master/contracts/RootChain/RootChain.sol#L165
    function startExit(
        bytes prevTx,
        bytes exitingTx, 
        bytes exitingTxSig
        /* bytes prevTxInclusionProof
        bytes exitingTxInclusionProof*/
    ) 
        external
    {
        ERC721PlasmaRLP.txData memory exitingTxData = exitingTx.createExitingTx(3);
        ERC721PlasmaRLP.txData memory prevTxData = prevTx.createExitingTx(3);
        uint prevBlock = exitingTxData.prevBlock;
        uint uid = exitingTxData.uid;
        require(uid == prevTxData.uid); // check that they are both for the same coin
        address owner = exitingTxData.owner;

        // Check if the exitingTx is signed by the previous tx's owner.
        bytes32 txHash = keccak256(exitingTx);
        require(
            txHash.ecverify(
                exitingTxSig, prevTxData.owner
            ), "Invalid sig"
        );

        // if prevTx is a deposit transaction, no need to provide a previous tx. Need to verify that the block root matches.
        if (prevBlock % childBlockInterval != 0 ) { 
            bytes32 prevTxHash = keccak256(prevTx);
            require(prevTxHash == childChain[prevBlock].root);
        } 
        // If exiting a transaction that was not in a deposit block, we need to check its merkle inclusion proof
        else {  // Block has to be submitted by operator that checkpoints any extra offchain transfers
             exitingTxSig;
        }

        uint priority = prevBlock * 10000000  + uid * 10000;
            // Also need to check signatures that match. It's OK if previous tx is invalid since someone will be able to challenge that exit as specified in the spec. 

        exitsQueue.insert(priority);
        exits[priority] = Exit({
            owner: owner, 
            tokenId: uid, 
            created_at: block.timestamp
        });
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

            // add token Id to the ones the owner can withdraw
            pendingWithdrawals[currentExit.owner].push(currentExit.tokenId);
            emit FinalizedExit(priority, currentExit.owner, currentExit.tokenId);
            emit CanWithdraw(currentExit.owner, currentExit.tokenId);

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

    // maybe add a count paramter to add a max amount of tokens the owner can withdraw
    // a pop function would be very helpful here for popping a certain number of tokens
    function withdraw() external {
        uint[] memory tokens = pendingWithdrawals[msg.sender];
        uint length = tokens.length;
        // for each token that is confirmed for withdrawal let the owner get it
        for (uint i = 0 ; i < length ; i ++ ) {
            cryptoCards.safeTransferFrom(address(this), msg.sender, tokens[i]);
        }
    }

    function getDepositBlock() public view returns (uint256) {
        return currentChildBlock.sub(childBlockInterval).add(currentDepositBlock);
    }

    /// receiver for erc721 to trigger a deposit
    function onERC721Received(address _from, uint256 _tokenId, bytes _data) 
        public 
        returns(bytes4) 
    {
        require(msg.sender == address(cryptoCards)); // can only be called by the associated cryptocards contract. 
        deposit(_from, _tokenId, _data);
        return ERC721_RECEIVED;
    }



}

