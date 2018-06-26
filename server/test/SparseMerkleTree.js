const utils = require('web3-utils');
const BN = require('bn.js');

module.exports = class SparseMerkleTree {
    constructor(depth, leaves) {
        this.depth = depth;
        // Initialize defaults
        this.defaultNodes = this.setdefaultNodes(depth);
        // Leaves must be a dictionary with key as the leaf's slot and value the leaf's hash
        this.leaves = leaves;

        if (leaves && Object.keys(leaves).length !== 0) {
            this.tree = this.createTree(this.leaves, this.depth, this.defaultNodes);
            this.root = this.tree[this.depth]['0'];
        } else {
            this.tree = [];
            this.root = this.defaultNodes[this.depth];
        }
    }

    setdefaultNodes(depth) {
        let defaultNodes = new Array(depth + 1);
        defaultNodes[0] = utils.soliditySha3(0);
        for (let i = 1; i < depth + 1; i++) {
            defaultNodes[i] = utils.soliditySha3(defaultNodes[i-1], defaultNodes[i-1]);
        }
        return defaultNodes;
    }

    createTree(orderedLeaves, depth, defaultNodes) {
        let tree = [orderedLeaves];
        let treeLevel = orderedLeaves;

        let nextLevel = {};
        let halfIndex;
        let value;

        for (let level = 0; level < depth; level++) {
            nextLevel = {};
            for(let index in treeLevel) {
                halfIndex = web3.toBigNumber(index).dividedToIntegerBy(2).toString();
                value = treeLevel[index];
                if (web3.toBigNumber(index).mod(2).isZero()) {
                    let coIndex = web3.toBigNumber(index).add(1).toString();
                    nextLevel[halfIndex] =
                        utils.soliditySha3(value, treeLevel[coIndex] || defaultNodes[level]);
                } else {
                    let coIndex = web3.toBigNumber(index).sub(1).toString();
                    if (treeLevel[coIndex] === undefined) {
                          nextLevel[halfIndex] =
                            utils.soliditySha3(defaultNodes[level], value);
                    }
                }
            }
            treeLevel = nextLevel;
            tree.push(treeLevel);
        }
        return tree;
    }

    createMerkleProof(uid) {
        let index = web3.toBigNumber(uid)
        let proof = '';
        let proofbits = new BN(0);
        let siblingIndex;
        let siblingHash;
        for (let level=0; level < this.depth; level++) {
            siblingIndex = index.mod(2).eq(0) ? index.add(1) : index.sub(1);
            index = index.dividedToIntegerBy(2);

            siblingHash = this.tree[level][siblingIndex.toString()];
            if (siblingHash) {
                proof += siblingHash.replace('0x', '');
                proofbits = proofbits.bincn(level);
            }
        }

        let buf = proofbits.toBuffer('be', 8);
        let total = Buffer.concat([buf, Buffer.from(proof, 'hex')]);
        return '0x' + total.toString('hex');
    }
}
