pragma solidity ^0.4.24;


// Based on https://rinkeby.etherscan.io/address/0x881544e0b2e02a79ad10b01eca51660889d5452b#code
// Original Code for the sparse merkle tree came from Wolkdb Plasma, this is now no longer compatible with that
// we have javascript and Golang implementations for reference of the new implementation.
contract SparseMerkleTree {

    uint8 constant DEPTH = 64;
    bytes32[DEPTH + 1] public defaultHashes;

    constructor() public {
        // defaultHash[0] is being set to keccak256(uint256(0));
        defaultHashes[0] = 0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563;
        setDefaultHashes(1, DEPTH);
    }

    function checkMembership(
        bytes32 leaf,
        bytes32 root,
        uint64 tokenID,
        bytes proof) public view returns (bool)
    {
        bytes32 computedHash = getRoot(leaf, tokenID, proof);
        return (computedHash == root);
    }

    // first 64 bits of the proof are the 0/1 bits
    function getRoot(bytes32 leaf, uint64 index, bytes proof) public view returns (bytes32) {
        require((proof.length - 8) % 32 == 0 && proof.length <= 2056);
        bytes32 proofElement;
        bytes32 computedHash = leaf;
        uint16 p = 8;
        uint64 proofBits;
        assembly {proofBits := div(mload(add(proof, 32)), exp(256, 24))}

        for (uint d = 0; d < DEPTH; d++ ) {
            if (proofBits % 2 == 0) { // check if last bit of proofBits is 0
                proofElement = defaultHashes[d];
            } else {
                p += 32;
                require(proof.length >= p);
                assembly { proofElement := mload(add(proof, p)) }
            }
            if (index % 2 == 0) {
                computedHash = keccak256(abi.encodePacked(computedHash, proofElement));
            } else {
                computedHash = keccak256(abi.encodePacked(proofElement, computedHash));
            }
            proofBits = proofBits / 2; // shift it right for next bit
            index = index / 2;
        }
        return computedHash;
    }

    function setDefaultHashes(uint8 startIndex, uint8 endIndex) private {
        for (uint8 i = startIndex; i <= endIndex; i ++) {
            defaultHashes[i] = keccak256(abi.encodePacked(defaultHashes[i-1], defaultHashes[i-1]));
        }
    }
}
