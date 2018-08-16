pragma solidity ^0.4.24;

import "./RLP.sol";


library Transaction {

    using RLP for bytes;
    using RLP for RLP.RLPItem;

    struct TX {
        uint64 slot;
        address owner;
        bytes32 hash;
        uint256 prevBlock;
        uint256 nonce;
        uint256 balance;
    }

    function getTx(bytes memory txBytes) internal pure returns (TX memory) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        TX memory transaction;

        transaction.slot = uint64(rlpTx[0].toUint());
        transaction.prevBlock = rlpTx[1].toUint();
        transaction.nonce = rlpTx[2].toUint();
        transaction.balance = rlpTx[3].toUint();
        transaction.owner = rlpTx[4].toAddress();
        if (transaction.prevBlock == 0) { // deposit transaction
            transaction.hash = keccak256(abi.encodePacked(transaction.slot));
        } else {
            transaction.hash = keccak256(txBytes);
        }
        return transaction;
    }

    function getHash(bytes memory txBytes) internal pure returns (bytes32 hash) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        uint64 slot = uint64(rlpTx[0].toUint());
        uint256 prevBlock = uint256(rlpTx[1].toUint());

        if (prevBlock == 0) { // deposit transaction
            hash = keccak256(abi.encodePacked(slot));
        } else {
            hash = keccak256(txBytes);
        }
    }

    function getNonce(bytes memory txBytes) internal pure returns (uint256 nonce) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        nonce = rlpTx[2].toUint();
    }

    function getBalance(bytes memory txBytes) internal pure returns (uint256 balance) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        balance = rlpTx[3].toUint();
    }

    function getOwner(bytes memory txBytes) internal pure returns (address owner) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(5);
        owner = rlpTx[4].toAddress();
    }
}
