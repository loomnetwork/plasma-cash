pragma solidity ^0.4.18;


library Merkle {
    function checkMembership(bytes32 leaf, uint256 index, bytes32 rootHash, bytes proof)
        internal
        pure
        returns (bool)
    {
        // Merkle tree depth is 32, each proof segment is 32 bytes
        // lg2(5) * 32 = 160
        // Needs to be improved with caching for empty trees. 
        require(proof.length == 160); 
        bytes32 proofElement;
        bytes32 computedHash = leaf;

        for (uint256 i = 32; i <= 160; i += 32) {
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
