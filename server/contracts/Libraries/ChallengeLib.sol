// Copyright Loom Network 2018 - All rights reserved, Dual licensed on GPLV3
// Learn more about Loom DappChains at https://loomx.io
// All derivitive works of this code must incluse this copyright header on every file 

pragma solidity ^0.4.24;

/**
* @title ChallengeLib
*
* ChallengeLib is a helper library for constructing challenges
*/

library ChallengeLib {
    struct Challenge {
        address owner;
        address challenger;
        bytes32 txHash;
        uint256 challengingBlockNumber;
    }

    function contains(Challenge[] storage _array, bytes32 txHash) internal view returns (bool) {
        int index = indexOf(_array, txHash);
        return index != -1;
    }

    function remove(Challenge[] storage _array, bytes32 txHash) internal returns (bool) {
        int index = indexOf(_array, txHash);
        if (index == -1) {
            return false; // Tx not in challenge arraey
        }
        // Replace element with last element
        Challenge memory lastChallenge = _array[_array.length - 1];
        _array[uint(index)] = lastChallenge;

        // Reduce array length
        delete _array[_array.length - 1];
        _array.length -= 1;
        return true;
    }

    function indexOf(Challenge[] storage _array, bytes32 txHash) internal view returns (int) {
        for (uint i = 0; i < _array.length; i++) {
            if (_array[i].txHash == txHash) {
                return int(i);
            }
        }
        return -1;
    }
}
