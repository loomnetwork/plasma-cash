# gocontracts

This repo is a bunch of libraries for Golang, to help emulate or speak with smart contracts written in solidity.



## Features:


## Solidity tightly packed formats

This library also supports creating Solidity's tightly packed data constructs, which are used together with sha3, sha256 and ripemd160 to create hashes.

```
Solidity code:

contract HashTest {
  function testSha3() returns (bytes32) {
   address addr1 = 0x43989fb883ba8111221e89123897538475893837;
   address addr2 = 0;
   uint val = 10000;
   uint timestamp = 1448075779;

   return sha3(addr1, addr2, val, timestamp); // will return 0xc3ab5ca31a013757f26a88561f0ff5057a97dfcc33f43d6b479abc3ac2d1d595
 }
}
Creating the same hash using this library:

import "github.com/loomnetwork/gocontracts"

gocontracts.SoliditySHA3( &[gocontracts.Pair]{
    { "address": "43989fb883ba8111221e89123897538475893837"}, "address", "uint", "uint" ],
    [ new BN(, 16), 0, 10000, 1448075779 ]
).toString('hex')
For the same data structure:

sha3 will return 0xc3ab5ca31a013757f26a88561f0ff5057a97dfcc33f43d6b479abc3ac2d1d595

```