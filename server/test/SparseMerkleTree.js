const utils = require('web3-utils');
// const BN = require('bn.js');
const BN = require('bignumber.js');

module.exports = class SparseMerkleTree {
  constructor(depth, leaves) {
    this.depth = depth;
    // Initialize defaults
    this.defaultNodes = this.setdefaultNodes(depth);
    this.leaves = leaves; // Leaves must be a dictionary with key as the leaf's slot and value the leaf's hash
    
    if (leaves) {
      this.tree = this.createTree(this.leaves, this.depth, this.defaultNodes)
      this.root = this.tree[this.depth-1][0]
    } else {
      this.tree = [];
      this.root = this.defaultNodes[this.depth-1];
    }
  }

  setdefaultNodes(depth) {
    let defaultNodes = new Array(depth);
    defaultNodes[0] = utils.soliditySha3(0);
    for (let i = 1; i < depth; i++) {
      defaultNodes[i] = utils.soliditySha3(defaultNodes[i-1], defaultNodes[i-1]);
    }
    return defaultNodes;
  }

  createTree(orderedLeaves, depth, defaultNodes) {
    let tree = [orderedLeaves];
    let treeLevel = orderedLeaves;

    let nextLevel = {};
    let prevIndex;
    let item;
    for (let level = 0; level < depth -1; level ++) {
      nextLevel = {};
      prevIndex = -1;
      for (var index in treeLevel) {
        var value = treeLevel[index]; 
        if (index % 2 === 0) {
          nextLevel[ Math.floor(index/2) ] = 
                  utils.soliditySha3(value + defaultNodes[level]);
        } else {
          if (index === prevIndex + 1) {
            nextLevel[Math.floor(index/2)] = utils.soliditySha3(value + defaultNodes[level]);
          } else {
            nextLevel[Math.floor(index/2)] = utils.soliditySha3(value + defaultNodes[level]);
          }
        }
      }

      // console.log(nextLevel);

      treeLevel = nextLevel;
      tree.push(treeLevel);
    }
    return tree;
  }

  createMerkleProof(uid) {
    let index = uid;
    let proof = '';
    let proofbits = '';
    let siblingIndex;
    let siblingHash;
    for (let level=0; level < this.depth -1; level++) {
      siblingIndex = index % 2 === 0 ? index + 1 : index -1;
      index = Math.floor(index / 2);

      siblingHash = this.tree[level][siblingIndex];
      if (siblingHash) {
        proof += siblingHash.slice(2, siblingHash.length)
        proofbits += '1'
      } else {
        proofbits += '0';
      }
    }
    let bits = new BN(proofbits, 2);
    console.log(proof);
    return utils.hexToBytes(utils.fromDecimal(bits));
  }

}
//// tx = (slot, prevblock, denom, newowner)
//slot = 1;
//prevblock = 1000;
//denom = 1;
//newowner = '0x123456789';
//data = [slot, prevblock, denom, newowner];
//tx = '0x' + RLP.encode(data).toString('hex'); 


// let leaves = {};
// leaves[slot] = utils.soliditySha3(tx);
// 
// slot = 5
// data = [slot, prevblock, denom, newowner];
// tx2 = '0x' + RLP.encode(data).toString('hex');
// leaves[slot] = utils.soliditySha3(tx2);
// 

// tree = new SparseMerkleTree(8, leaves);
// 
// tree.createMerkleProof(1);
// 
// tree = new SparseMerkleTree(4);
