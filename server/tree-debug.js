SparseMerkleTree = require('./test/SparseMerkleTree.js');
utils = require('web3-utils');

let slot = 2
let txHash = '0xcf04ea8bb4ff94066eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84'
let leaves = {};
leaves[slot] = txHash;

//let slot2 = 100
//let txHash2 = '0xcf04aaaaaaaaaaaa6eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84'
//leaves[slot2] = txHash2;

let tree = new SparseMerkleTree(64, leaves);

// console.log(tree)

console.log('ROOT:', tree.root);
console.log('Proof', tree.createMerkleProof(slot)) // all \x00 so not printable

