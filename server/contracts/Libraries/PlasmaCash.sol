pragma solidity ^0.4.23;

improt './TxHash.sol';

contract PlasmaCash is sparseMerkle {

    event Deposit(address _depositor, uint64 _depositIndex, uint64 _denomination, uint64 _tokenID);
    event DepositExit(address _exiter, uint64 _depositIndex, uint64 _denomination, uint64 _tokenID, uint256 _timestamp);
    event StartExit(address _exiter, uint64 _depositIndex, uint64 _denomination, uint64 _tokenID, uint256 _timestamp);
    event PublishedBlock(bytes32 _rootHash, uint64 _blknum, uint64 _currentDepositIndex);


    //using RLP for bytes;
    //using RLP for RLP.RLPItem;
    //using RLP for RLP.Iterator;

    using TxHash for bytes;
    address public authority;
    uint64 public currentDepositIndex;
    uint64 public currentBlkNum;
    mapping(uint64 => bytes32) public childChain;
    mapping(uint64 => uint64) public depositIndex;
    mapping(uint64 => uint64) public depositBalance;


    modifier isAuthority() {
        require(msg.sender == authority);
        _;
    }


   constructor() public {
        authority = msg.sender;
        currentDepositIndex = 0;
        currentBlkNum = 0;
    }

    // @dev Allows Plasma chain operator to submit block root
    // @param blkRoot The root of a child chain block
    // @param blknum The child chain block number
    function submitBlock(bytes32 _blkRoot, uint64 _blknum) public isAuthority {
        //require(currentBlkNum + 1 == _blknum);
        //currentBlkNum += 1;
        childChain[_blknum] = _blkRoot;
        currentBlkNum = _blknum;
        emit PublishedBlock(_blkRoot, _blknum, currentDepositIndex);
    }

    // @dev Allows anyone to deposit eth into the Plasma chain, Reject tokendeposit for now.
    function deposit() public payable {
        require (msg.value < (2**64 - 2) ); //18.446744073709551615Eth
        uint64 depositAmount = uint64(msg.value % (2 ** 64)) ;
        uint64 tokenID = uint64(uint256(keccak256(msg.sender, currentDepositIndex, depositAmount)) % (2 ** 64));
        require (depositAmount > 0 && depositBalance[tokenID] == 0);
        depositBalance[tokenID] = depositAmount;
        depositIndex[currentDepositIndex] = tokenID;
        emit Deposit(msg.sender, currentDepositIndex, depositAmount, tokenID);
        currentDepositIndex += 1;
    }

    // @dev Allows original owner to submit withdraw request [bond will be added in futre]
    function depositExit(uint64 _depositIndex) public {
        uint64 tokenID = depositIndex[_depositIndex];
        uint64 denomination = depositBalance[tokenID];
        require(uint64(uint256(keccak256(msg.sender, _depositIndex, denomination)) % (2 ** 64)) == tokenID);
        require(denomination > 0);
        depositBalance[tokenID] = 0;
        //TODO: Adding depositExit to PriorityQueue
        msg.sender.transfer(denomination);
        emit DepositExit(msg.sender, _depositIndex, denomination, tokenID, block.timestamp);
    }

    // @ dev Takes in the transaction transfering ownership to the current owner and the proofs necesary to prove there inclusion
    function startExit(uint64 tokenID, bytes txBytes1, bytes txBytes2, bytes proof1, bytes proof2, uint64 blk1, uint64 blk2) public {

        TxHash.TX memory tx2 = txBytes2.getTx();
        TxHash.TX memory tx1 = txBytes1.getTx();
        require(tx2.Recipient == msg.sender, "unauth exit");
        require(tx1.TokenId == tokenID && tx2.TokenId == tokenID, "tokenID mismatch");
        require(txBytes1.verifyTX(), "tx1 sig failure");
        require(txBytes2.verifyTX(), "tx2 sig failure");

        //checkMembership(leaf, root, tokenID, proof);
        require(sparseMerkle.checkMembership(keccak256(txBytes1),childChain[blk1],tx1.TokenId, proof1), "tx1 non member");
        require(sparseMerkle.checkMembership(keccak256(txBytes2),childChain[blk2],tx2.TokenId, proof2), "tx2 non member");

        //TODO: PriorityQueue: addExitToQueue(exitUid, exitor, amount);
        emit StartExit(msg.sender, tx2.DepositIndex, tx2.Denomination, tokenID, block.timestamp);

    }

    function () public payable {
        deposit();
    }

    // @dev Submit proof of a transaction before prevTx forcing the exitor to call respondChallengeBefore and submit a transaction before prevTx or prevTx itself.
    function challengeBefore(uint exitUid, bytes txBytes, bytes proof) public {

    }

    // @dev Submit a transaction after the challengeBefore tx and before prevTx or prevTx itself.
    function respondChallengeBefore(uint exitUid, bytes txBytes, bytes proof) public {

    }

    // @ dev Submit proof that the exiting uid has been spent between the prevTx and tx, proving that the owner of the exiting uid is illegitimate.
    function challengeBetween(uint exitUid, bytes txBytes, bytes proof) public {

    }

    // @ dev Submit proof that the exiting uid has been spent on the child chain after the exit has been triggered.
    function challengeAfter(uint exitUid, bytes txBytes, bytes proof) public {

    }

    // @dev Finalizes an exit after its challenge period has finished, removing uid from wallet if amountLeft is 0 or updating it in the case of a split coin exit.
    function finalizeExit(uint exitUid) public {

    }

}
