pragma solidity ^0.4.24;

/**
* @title Challenge
*
* Challenge is used to constructe Challenge.
*/

library Challenges {
    struct Challenge {
        bool hasValue;
        bytes ChallengeTx;
        uint ChallengeTxBlkNum;
    }

    function contains(Challenge[] storage _array, bytes _ChallengeTx) internal view returns (bool) {
        int index = _indexOf(_array, _ChallengeTx);
        return index != -1;
    }

    function remove(Challenge[] storage _array, bytes _ChallengeTx) internal returns (bool) {
        int index = _indexOf(_array, _ChallengeTx);
        if (index == -1) {
            return false;
        }
        Challenge memory lastChallenge = _array[_array.length - 1];
        _array[uint(index)] = lastChallenge;

        delete _array[_array.length - 1].hasValue;
        _array.length -= 1;
        return true;
    }

    function _indexOf(Challenge[] storage _array, bytes _ChallengeTx) private view returns (int) {
        for (uint i = 0; i < _array.length; i++) {
            bytes memory a = _array[i].ChallengeTx;
            bytes memory b = _ChallengeTx;

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
