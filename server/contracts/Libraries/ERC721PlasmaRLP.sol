pragma solidity ^0.4.23;

/**
 * @title RLPReader
 * @dev RLPReader is used to read and parse RLP encoded data in memory.
 * @author Andreas Olofsson (androlo1980@gmail.com)
 */

library ERC721PlasmaRLP {
    uint constant DATA_SHORT_START = 0x80;
    uint constant DATA_LONG_START = 0xB8;
    uint constant LIST_SHORT_START = 0xC0;
    uint constant LIST_LONG_START = 0xF8;

    uint constant DATA_LONG_OFFSET = 0xB7;


    struct RLPItem {
        uint _unsafe_memPtr;    // Pointer to the RLP-encoded bytes.
        uint _unsafe_length;    // Number of bytes. This is the full length of the string.
    }

    struct txData {
        uint256 slot;
        uint256 prevBlock;
        uint256 denomination;
        address owner;
        bytes sig;
    }

    struct Iterator {
        RLPItem _unsafe_item;   // Item that's being iterated over.
        uint _unsafe_nextPtr;   // Position of the next item in the list.
    }

    /* Public Functions */

    function getTxData(bytes memory txDataBytes)
        internal
        pure
        returns (txData)
    {
        RLPItem[] memory txList = toList(toRLPItem(txDataBytes), 5);
        return txData({
            slot: toUint(txList[0]),
            prevBlock: toUint(txList[1]),
            denomination: toUint(txList[2]),
            owner: toAddress(txList[3]),
            sig: toBytes(txList[4])
        });
    }

    /* Iterator */

    function next(Iterator memory self) private pure returns (RLPItem memory subItem) {
        uint ptr = self._unsafe_nextPtr;
        uint itemLength = _itemLength(ptr);
        subItem._unsafe_memPtr = ptr;
        subItem._unsafe_length = itemLength;
        self._unsafe_nextPtr = ptr + itemLength;
    }

    function hasNext(Iterator memory self) private pure returns (bool) {
        RLPItem memory item = self._unsafe_item;
        return self._unsafe_nextPtr < item._unsafe_memPtr + item._unsafe_length;
    }

    /* RLPItem */

    /// @dev Creates an RLPItem from an array of RLP encoded bytes.
    /// @param self The RLP encoded bytes.
    /// @return An RLPItem
    function toRLPItem(bytes memory self) private pure returns (RLPItem memory) {
        uint len = self.length;
        uint memPtr;
        assembly {
            memPtr := add(self, 0x20)
        }
        return RLPItem(memPtr, len);
    }

    /// @dev Create an iterator.
    /// @param self The RLP item.
    /// @return An 'Iterator' over the item.
    function iterator(RLPItem memory self) private pure returns (Iterator memory it) {
        uint ptr = self._unsafe_memPtr + _payloadOffset(self);
        it._unsafe_item = self;
        it._unsafe_nextPtr = ptr;
    }

    /// @dev Get the list of sub-items from an RLP encoded list.
    /// Warning: This requires passing in the number of items.
    /// @param self The RLP item.
    /// @return Array of RLPItems.
    function toList(RLPItem memory self, uint256 numItems) private pure returns (RLPItem[] memory list) {
        list = new RLPItem[](numItems);
        Iterator memory it = iterator(self);
        uint idx;
        while(idx < numItems) {
            list[idx] = next(it);
            idx++;
        }
    }

    /// @dev Decode an RLPItem into a uint. This will not work if the
    /// RLPItem is a list.
    /// @param self The RLPItem.
    /// @return The decoded string.
    function toUint(RLPItem memory self) private pure returns (uint data) {
        uint rStartPos; uint len;
        (rStartPos, len) = _decode(self);
        assembly {
            data := div(mload(rStartPos), exp(256, sub(32, len)))
        }
    }

    /// @dev Decode an RLPItem into an address. This will not work if the
    /// RLPItem is a list.
    /// @param self The RLPItem.
    /// @return The decoded string.
    function toAddress(RLPItem memory self)
        private
        pure
        returns (address data)
    {
        uint rStartPos; uint len;
        (rStartPos, len) = _decode(self);
        assembly {
            data := div(mload(rStartPos), exp(256, 12))
        }
    }

    // Get the payload offset.
    function _payloadOffset(RLPItem memory self)
        private
        pure
        returns (uint)
    {
        uint b0;
        uint memPtr = self._unsafe_memPtr;
        assembly {
            b0 := byte(0, mload(memPtr))
        }
        if(b0 < DATA_SHORT_START)
            return 0;
        if(b0 < DATA_LONG_START || (b0 >= LIST_SHORT_START && b0 < LIST_LONG_START))
            return 1;
    }

    // Get the full length of an RLP item.
    function _itemLength(uint memPtr)
        private
        pure
        returns (uint len)
    {
        uint b0;
        assembly {
            b0 := byte(0, mload(memPtr))
        }
        if (b0 < DATA_SHORT_START)
            len = 1;
        else if (b0 < DATA_LONG_START)
            len = b0 - DATA_SHORT_START + 1;
    }

    // Get start position and length of the data.
    function _decode(RLPItem memory self)
        private
        pure
        returns (uint memPtr, uint len)
    {
        uint b0;
        uint start = self._unsafe_memPtr;
        assembly {
            b0 := byte(0, mload(start))
        }
        if (b0 < DATA_SHORT_START) {
            memPtr = start;
            len = 1;
            return;
        }
        if (b0 < DATA_LONG_START) {
            len = self._unsafe_length - 1;
            memPtr = start + 1;
        } else {
            uint bLen;
            assembly {
                bLen := sub(b0, 0xB7) // DATA_LONG_OFFSET
            }
            len = self._unsafe_length - 1 - bLen;
            memPtr = start + bLen + 1;
        }
        return;
    }
    
    /// @dev Return the RLP encoded bytes.
    /// @param self The RLPItem.
    /// @return The bytes.
    function toBytes(RLPItem memory self)
        internal
        pure
        returns (bytes memory bts)
    {
        uint len = self._unsafe_length;
        if (len == 0)
            return;
        bts = new bytes(len);
        _copyToBytes(self._unsafe_memPtr, bts, len);
    }

    /// @dev Decode an RLPItem into bytes. This will not work if the
    /// RLPItem is a list.
    /// @param self The RLPItem.
    /// @return The decoded string.
    function toData(RLPItem memory self)
        internal
        pure
        returns (bytes memory bts)
    {
        // require(isData(self));
        uint rStartPos; uint len;
        (rStartPos, len) = _decode(self);
        bts = new bytes(len);
        _copyToBytes(rStartPos, bts, len);
    }

    // Assumes that enough memory has been allocated to store in target.
    function _copyToBytes(uint btsPtr, bytes memory tgt, uint btsLen)
        private
        pure
    {
        // Exploiting the fact that 'tgt' was the last thing to be allocated,
        // we can write entire words, and just overwrite any excess.
        assembly {
            {
                // evm operations on words
                let words := div(add(btsLen, 31), 32)
                let rOffset := btsPtr
                let wOffset := add(tgt, 0x20)
                for
                    { let i := 0 } // start at arr + 0x20 -> first byte corresponds to length
                    lt(i, words)
                    { i := add(i, 1) }
                {
                    let offset := mul(i, 0x20)
                    mstore(add(wOffset, offset), mload(add(rOffset, offset)))
                }
                mstore(add(tgt, add(0x20, mload(tgt))), 0)
            }
        }

	}

}
