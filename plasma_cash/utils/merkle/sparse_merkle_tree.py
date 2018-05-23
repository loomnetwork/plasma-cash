from collections import OrderedDict

# from ethereum.utils import w3.sha3
from web3.auto import w3
from hexbytes import HexBytes


class SparseMerkleTree(object):

    def __init__(self, depth=64, leaves={}):
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
        default_hash = w3.soliditySha3(['bytes32'], [HexBytes('00' * 32 )])
        default_nodes = [default_hash]
        for level in range(1, depth):
            prev_default = default_nodes[level - 1]
            default_nodes.append(w3.soliditySha3(['bytes32'], [prev_default + prev_default]))
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
                    next_level[index // 2] = w3.soliditySha3(['bytes32'], [value + default_nodes[level]])
                else:
                    # If the node is a right node, check if its left sibling is a default node.
                    if index == prev_index + 1:
                        next_level[index // 2] = w3.soliditySha3(['bytes32'], [tree_level[prev_index] + value])
                    else:
                        next_level[index // 2] = w3.soliditySha3(['bytes32'], [default_nodes[level] + value])
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
        proofbits = b''
        for level in range(self.depth - 1):
            sibling_index = index + 1 if index % 2 == 0 else index - 1
            index = index // 2
            if sibling_index in self.tree[level]:
                proof += self.tree[level][sibling_index]
                proofbits += b'1'

            else:
                proofbits += b'0'

        # Need to convert the binary string to bytes for solidity to understand.
        proof_bytes = self._to_bytes(proofbits)
        return proof_bytes + proof

    def _to_bytes(self, bits):
        # lord have mercy on my soul for this hack. 
        # There is probably a better way to store the slots that were 0/1 in a number
        # struct.pack('>Q', 0x200000000000000) <- hexstr does the job
        b = w3.toBytes(hexstr=hex(int(bits,2)))
        b  = b.rjust(self.depth // 8, b'\0') # pad to bytes needed
        return b

    class TreeSizeExceededException(Exception):
        """there are too many leaves for the tree to build"""
