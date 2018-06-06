pragma solidity ^0.4.24;

import "./RLP.sol";


library Transaction {

    using RLP for bytes;
    using RLP for RLP.RLPItem;

    struct TX {
        uint64 slot;
        uint32 denomination; // 2**32 more than enough for NFTs, helps pack in 1 slot
        address owner;
        uint256 prevBlock;
    }

    function getTx(bytes memory txBytes) internal pure returns (TX memory) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(4);
        TX memory transaction;

        transaction.slot = uint64(rlpTx[0].toUint());
        transaction.prevBlock = uint64(rlpTx[1].toUint());
        transaction.denomination = uint32(rlpTx[2].toUint());
        transaction.owner = rlpTx[3].toAddress();

        return transaction;
    }

}
