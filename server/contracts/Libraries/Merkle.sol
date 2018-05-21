pragma solidity ^0.4.18;


library Merkle {
    uint constant PROOF_LENGTH = 992;

    function checkMembership(bytes32 leaf, uint256 index, bytes32 rootHash, bytes proof)
        internal
        pure
        returns (bool)
    {
        // Merkle tree depth is 32, each proof segment is 32 bytes
        // Need 31 items in the proof so total length of proof is 31 * 32 = 992
        // Can be further improved for SMT
        require(proof.length == PROOF_LENGTH); 
        bytes32 proofElement;
        bytes32 computedHash = leaf;

        for (uint256 i = 32; i <= PROOF_LENGTH; i += 32) {
            assembly {
                proofElement := mload(add(proof, i))
            }
            if (index % 2 == 0) {
                computedHash = keccak256(computedHash, proofElement);
            } else {
                computedHash = keccak256(proofElement, computedHash);
            }
            index = index / 2;
        }
        return computedHash == rootHash;
    }
}
