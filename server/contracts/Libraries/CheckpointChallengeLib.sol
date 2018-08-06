pragma solidity ^0.4.24;

/**
* @title Challenge
*
* Challenge is used to construct challenges regarding checkpoints.
*/

library CheckpointChallengeLib {
    struct CheckpointChallenge {
        bool hasValue;
        bytes CheckpointChallengeTx;
        uint CheckpointChallengeTxBlkNum;
    }

    function contains(CheckpointChallenge[] storage _array, bytes _CheckpointChallengeTx) internal returns (bool) {
        int index = _indexOf(_array, _CheckpointChallengeTx);
        return index != -1;
    }

    function remove(CheckpointChallenge[] storage _array, bytes _CheckpointChallengeTx) internal returns (bool) {
        int index = _indexOf(_array, _CheckpointChallengeTx);
        if (index == -1) {
            return false;
        }
        CheckpointChallenge memory lastChallenge = _array[_array.length - 1];
        _array[uint(index)] = lastChallenge;

        delete _array[_array.length - 1].hasValue;
        _array.length -= 1;
        return true;
    }

    function _indexOf(CheckpointChallenge[] storage _array, bytes _CheckpointChallengeTx) private returns (int) {
        for (uint i = 0; i < _array.length; i++) {
            bytes memory a = _array[i].CheckpointChallengeTx;
            bytes memory b = _CheckpointChallengeTx;

            if (compare(a, b) == 0) {
                return int(i);
            }
        }
        return -1;
    }

    function compare(bytes _a, bytes _b) private pure returns (int) {
        uint minLength = _a.length;
        if (_b.length < minLength) {
            minLength = _b.length;
        }

        for (uint i = 0; i < minLength; i++) {
            if (_a[i] < _b[i]) {
                return -1;
            } else if (_a[i] > _b[i]) {
                return 1;
            }
        }
        if (_a.length < _b.length) {
            return -1;
        } else if (_a.length > _b.length) {
            return 1;
        } else {
            return 0;
        }
    }
}
