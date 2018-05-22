pragma solidity ^0.4.23;


import './RLP.sol';
import './RLPEncode.sol';
import './ECVerify.sol';

library TxHash {
    
    using RLP for bytes;
    using RLP for RLP.RLPItem;
    using RLP for RLP.Iterator;
    using RLPEncode for bytes[];
    using RLPEncode for bytes;
    
    struct TX {
        uint64 TokenId;
        uint64 Denomination;
        uint64 DepositIndex;
        uint64 PrevBlock;
        address PrevOwner;
        address Recipient;
        uint8   PTransfer;
        uint8   Receipt;
    }
    
    struct RLPItem {
        uint _unsafe_memPtr;    // Pointer to the RLP-encoded bytes.
        uint _unsafe_length;    // Number of bytes. This is the full length of the string.
    }


    function constructUnsignedHash(bytes memory txBytes) internal view returns (bytes32) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        bytes[] memory unsignedTx = new bytes[](9);
        for(uint i=0; i<rlpTx.length; i++) {
            if (i < 7){
                unsignedTx[i] = rlpTx[i].toBytes();
            }else{
                unsignedTx[i] = new bytes(0).encodeBytes();
            }
        }
        bytes memory rlpUnsignedTx = unsignedTx.encodeList();
        return keccak256(rlpUnsignedTx);
    }

    function constructUnsigned(bytes memory txBytes) internal view returns (bytes) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        bytes[] memory unsignedTx = new bytes[](9);
        for(uint i=0; i<rlpTx.length; i++) {
            if (i < 7){
                unsignedTx[i] = rlpTx[i].toBytes();
            }else{
                unsignedTx[i] = new bytes(0).encodeBytes();
            }
        }
        return unsignedTx.encodeList();
    }
     
    function getSig(bytes memory txBytes) internal view returns (bytes) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        return rlpTx[7].toData();
     }

    function getSigner(bytes memory txBytes) internal view returns (address signer) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        bytes[] memory unsignedTx = new bytes[](9);
        bytes memory sig;
        address prevOwner;
            
        for(uint i=0; i<rlpTx.length; i++) {
            if (i == 4){
                prevOwner = rlpTx[i].toAddress();
                unsignedTx[i] = rlpTx[i].toBytes();
            }else if (i == 7){
                sig = rlpTx[i].toData();
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else if (i == 8){
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else {
                unsignedTx[i] = rlpTx[i].toBytes();
            }
        }
        bytes memory rlpUnsignedTx =  unsignedTx.encodeList();
        bytes32 txHash = keccak256(rlpUnsignedTx);
        return txHash.recover(sig);
     }
     
    function verifyTX(bytes memory txBytes) internal view returns (bool) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        if ((rlpTx[3].toUint() == 0)) {
            //If prevBlock = 0, check whether last 8 bytes of keccak256(recipient, depositIndex, denomination) == tokenID
            return uint256(
                keccak256(
                    rlpTx[5].toAddress(),
                    uint64(rlpTx[2].toUint()),
                    uint64(rlpTx[1].toUint()))
                ) 
                % (2**64) 
                == (rlpTx[0].toUint());
        }
        bytes[] memory unsignedTx = new bytes[](9);
        bytes memory sig;
        address prevOwner;
        for(uint i=0; i<9; i++) {
            if (i == 4){
                prevOwner = rlpTx[i].toAddress();
                unsignedTx[i] = rlpTx[i].toBytes();                
            }else if (i == 7){
                sig = rlpTx[i].toData();
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else if (i == 8){
                unsignedTx[i] = new bytes(0).encodeBytes();
            }else {
                unsignedTx[i] = rlpTx[i].toBytes();
            }
        }
        bytes memory rlpUnsignedTx = unsignedTx.encodeList();
        bytes32 txHash = keccak256(rlpUnsignedTx);
        return txHash.ecverify(sig, prevOwner);
     }
     

    function getPrevOwner(bytes memory txBytes) internal view returns (address) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        return rlpTx[4].toAddress();
    }
     
    function getTokenID(bytes memory txBytes) internal view returns (uint64) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        return uint64(rlpTx[0].toUint());
    }
     
    function getTx(bytes memory txBytes) internal view returns (TX memory) {
        RLP.RLPItem[] memory rlpTx = txBytes.toRLPItem().toList(9);
        TX memory tx;
        tx.TokenId =  uint64(rlpTx[0].toUint());
        tx.Denomination =  uint64(rlpTx[1].toUint());
        tx.DepositIndex =  uint64(rlpTx[2].toUint());
        tx.PrevBlock =  uint64(rlpTx[3].toUint());
        tx.PrevOwner =  rlpTx[4].toAddress();
        tx.Recipient =  rlpTx[5].toAddress();
        tx.PTransfer =  uint8(rlpTx[6].toUint());
        tx.Receipt =  uint8(rlpTx[8].toUint());
        return tx;
    }     
}
