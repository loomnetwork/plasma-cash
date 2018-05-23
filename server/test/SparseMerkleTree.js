const utils = require('web3-utils');
const BN = require('bn.js');

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
                  utils.soliditySha3(value, defaultNodes[level]);
        } else {
          if (index === prevIndex + 1) {
            nextLevel[Math.floor(index/2)] = utils.soliditySha3(treeLevel[prevIndex], value);
          } else {
            nextLevel[Math.floor(index/2)] = utils.soliditySha3(defaultNodes[level], value);
          }
        }
        prevIndex = index;
      }
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
        proof += siblingHash.replace('0x', '')
        proofbits += '1'
      } else {
        proofbits += '0';
      }
    }

    let reversed = proofbits.split("").reverse().join("");
    let bits = new BN(reversed, 2);
    // Must convert the BN to '\x12\x34' hexstring. Currently buggy. Uncomment the second transaction in testPlasma.js to test this
    let buf = bits.toBuffer().toString().padStart(8, '\x00')
    return buf + proof;
  }

}
