pragma solidity ^0.4.21;

// This contract exposes the precompiled methods built into Loom binary
// Don't use this class directly

contract DappChain {

    function loomtransferUxto(bytes data) public view {
        callPFAssembly(9, data);
    }

 
    function callPF(uint32 _addr, bytes _input) public view returns (bool) {
        address addr = _addr;
        return addr.call(_input);
    }

    uint constant ArraySize = 10;
    function callPFAssembly(uint64 _addr, bytes _input) public view returns (uint256[ArraySize]) {
        address addr = _addr;
        uint256 inSize = _input.length * 4 + 1;
        uint256 outSize = ArraySize * 0x20;
        uint256[ArraySize] memory rtv;
        assembly{
            let start := add(_input, 0x04)
            if iszero(call(
                5000,
                addr,
                0,
                start,
                inSize,
                rtv,
                outSize
            )) {
                revert(0,0)
            }
            mstore(0x40, add(0x40, outSize))
        }
        return rtv;
    }

}