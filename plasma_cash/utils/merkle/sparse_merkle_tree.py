from collections import OrderedDict

# from ethereum.utils import w3.sha3
from web3.auto import w3
from hexbytes import HexBytes


class SparseMerkleTree(object):

    def __init__(self, depth=256, leaves={}):
        self.depth = depth
        if len(leaves) > 2**(depth-1):
            raise self.TreeSizeExceededException(
                'tree with depth {} could not have {} leaves'.format(depth, len(leaves))
            )

        # Sort the transaction dict by index.
        self.leaves = OrderedDict(sorted(leaves.items(), key=lambda t: t[0]))
        self.default_nodes = self.create_default_nodes(self.depth)
        if leaves:
            self.tree = self.create_tree(self.leaves, self.depth, self.default_nodes)
            self.root = self.tree[-1][0]
        else:
            self.tree = []
            self.root = self.default_nodes[self.depth - 1]

    def create_default_nodes(self, depth):
        # Default nodes are the nodes whose children are both empty nodes at each level.
        default_hash = w3.sha3(HexBytes('00' * 32 ))
        default_nodes = [default_hash]
        for level in range(1, depth):
            prev_default = default_nodes[level - 1]
            default_nodes.append(w3.sha3(prev_default + prev_default))
        return default_nodes

    def create_tree(self, ordered_leaves, depth, default_nodes):
        tree = [ordered_leaves]
        tree_level = ordered_leaves
        for level in range(depth - 1):
            next_level = {}
            prev_index = -1
            for index, value in tree_level.items():
                if index % 2 == 0:
                    # If the node is a left node, assume the right sibling is a default node.
                    # in the case right sibling is not default node,
                    # it would override on next round
                    next_level[index // 2] = w3.sha3(value + default_nodes[level])
                else:
                    # If the node is a right node, check if its left sibling is a default node.
                    if index == prev_index + 1:
                        next_level[index // 2] = w3.sha3(tree_level[prev_index] + value)
                    else:
                        next_level[index // 2] = w3.sha3(default_nodes[level] + value)
                prev_index = index
            tree_level = next_level
            tree.append(tree_level)
        return tree

    def create_merkle_proof(self, uid):
        # Generate a merkle proof for a leaf with provided index. A proof is the concatenation of
        # the hash of node's sibling from leaf to root.
        index = uid
        proof = b''
        for level in range(self.depth - 1):
            sibling_index = index + 1 if index % 2 == 0 else index - 1
            index = index // 2
            if sibling_index in self.tree[level]:
                proof += self.tree[level][sibling_index]
            else:
                proof += self.default_nodes[level]
        return proof

    class TreeSizeExceededException(Exception):
        """there are too many leaves for the tree to build"""
