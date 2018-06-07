from collections import OrderedDict
from eth_utils.crypto import keccak


class SparseMerkleTree(object):

    def __init__(self, depth=64, leaves={}):
        self.depth = depth
        if len(leaves) > 2**(depth-1):
            raise self.TreeSizeExceededException(
                'tree with depth {} cannot have {} leaves'.format(depth,
                                                                  len(leaves)))

        # Sort the transaction dict by index.
        self.leaves = OrderedDict(sorted(leaves.items(), key=lambda t: t[0]))
        self.default_nodes = self.create_default_nodes(self.depth)
        if leaves:
            self.tree = self.create_tree(self.leaves, self.depth,
                                         self.default_nodes)
            self.root = self.tree[-1][0]
        else:
            self.tree = []
            self.root = self.default_nodes[self.depth - 1]

    def create_default_nodes(self, depth):
        # Default nodes are the nodes whose children are both empty nodes at
        # each level.
        default_hash = keccak(b'\x00' * 32)
        default_nodes = [default_hash]
        for level in range(1, depth):
            prev_default = default_nodes[level - 1]
            default_nodes.append(
                keccak(prev_default * 2))
        return default_nodes

    def create_tree(self, ordered_leaves, depth, default_nodes):
        tree = [ordered_leaves]
        tree_level = ordered_leaves
        for level in range(depth - 1):
            next_level = {}
            prev_index = -1
            for index, value in tree_level.items():
                if index % 2 == 0:
                    # If the node is a left node, assume the right sibling is
                    # a default node. In the case right sibling is not default
                    # node, it would override on next round
                    next_level[index // 2] = keccak(value
                                                    + default_nodes[level])
                else:
                    # If the node is a right node, check if its left sibling is
                    # a default node.
                    if index == prev_index + 1:
                        next_level[index // 2] = keccak(tree_level[prev_index]
                                                        + value)
                    else:
                        next_level[index // 2] = keccak(default_nodes[level]
                                                        + value)
                prev_index = index
            tree_level = next_level
            tree.append(tree_level)
        return tree

    def create_merkle_proof(self, uid):
        # Generate a merkle proof for a leaf with provided index.
        # First `depth/8` bytes of the proof are necessary for checking if
        # we are at a default-node
        index = uid
        proof = b''
        proofbits = 0
        for level in range(self.depth - 1):
            sibling_index = index + 1 if index % 2 == 0 else index - 1
            index = index // 2
            if sibling_index in self.tree[level]:
                proof += self.tree[level][sibling_index]
                proofbits += 2 ** level

        proof_bytes = proofbits.to_bytes(8, byteorder='big')
        return proof_bytes + proof

    def verify(self, uid, proof):
        ''' Checks if the proof for the leaf at `uid` is valid'''
        # assert (len(proof) -8 % 32) == 0
        assert len(proof) <= 2056

        proofbits = int.from_bytes((proof[0:8]), byteorder='big')
        index = uid
        p = 8
        computed_hash = self.leaves[index]

        for d in range(self.depth-1):
            if (proofbits % 2 == 0):
                proof_element = self.default_nodes[d]
            else:
                proof_element = proof[p:p+32]
                p += 32
            if (index % 2 == 0):
                computed_hash = keccak(computed_hash + proof_element)
            else:
                computed_hash = keccak(proof_element + computed_hash)
            proofbits = proofbits // 2
            index = index // 2
        return computed_hash == self.root

    class TreeSizeExceededException(Exception):
        """there are too many leaves for the tree to build"""
