import pytest
from hexbytes import HexBytes
from eth_utils.crypto import keccak
from utils.merkle.sparse_merkle_tree import SparseMerkleTree

empty_val = b'\x00' * 32
default_hash = keccak(empty_val)

class TestSparseMerkleTree(object):

    def test_size_limits(self):
        with pytest.raises(SparseMerkleTree.TreeSizeExceededException):
            SparseMerkleTree(depth=1, leaves={0: empty_val, 1: empty_val})

    def test_empty_SMT(self):
        emptyTree = SparseMerkleTree(64, {})
        assert len(emptyTree.leaves) == 0

    def test_all_leaves_with_val(self):
        dummy_val = b'\x01' * 32
        leaves = {0: dummy_val, 1: dummy_val, 2: dummy_val, 3: dummy_val}
        tree = SparseMerkleTree(depth=3, leaves=leaves)
        mid_level_val = keccak(dummy_val + dummy_val)
        assert tree.root == keccak(mid_level_val + mid_level_val)

    def test_empty_leaves(self):
        tree = SparseMerkleTree(depth=3)
        mid_level_val = keccak(default_hash * 2)
        assert tree.root == keccak(mid_level_val * 2)

    def test_empty_left_leave(self):
        dummy_val = b'\x01' * 32
        leaves = {1: dummy_val, 2: dummy_val, 3: dummy_val}
        tree = SparseMerkleTree(depth=3, leaves=leaves)
        mid_left_val = keccak(default_hash + dummy_val)
        mid_right_val = keccak(dummy_val + dummy_val)
        assert tree.root == keccak(mid_left_val + mid_right_val)

    def test_empty_right_leave(self):
        dummy_val = b'\x01' * 32
        leaves = {0: dummy_val, 2: dummy_val, 3: dummy_val}
        tree = SparseMerkleTree(depth=3, leaves=leaves)
        mid_left_val = keccak(dummy_val + default_hash)
        mid_right_val = keccak(dummy_val + dummy_val)
        assert tree.root == keccak(mid_left_val + mid_right_val)

    def test_exceed_tree_size(self):
        with pytest.raises(SparseMerkleTree.TreeSizeExceededException):
            SparseMerkleTree(depth=1, leaves={0: '0', 1: '1'})

    def test_create_merkle_proof(self):
        dummy_val = keccak(2)
        dummy_val_2 = keccak(3)
        leaves = {0: dummy_val, 2: dummy_val, 3: dummy_val_2}
        tree = SparseMerkleTree(depth=3, leaves=leaves)
        mid_left_val = keccak(dummy_val + default_hash)
        mid_right_val = keccak(dummy_val + dummy_val_2)
        assert tree.create_merkle_proof(0) == (2).to_bytes(8, byteorder='big') + mid_right_val
        assert tree.create_merkle_proof(1) == (3).to_bytes(8, byteorder='big') + dummy_val + mid_right_val
        assert tree.create_merkle_proof(2) == (3).to_bytes(8, byteorder='big') + dummy_val_2 + mid_left_val
        assert tree.create_merkle_proof(3) == (3).to_bytes(8, byteorder='big') + dummy_val + mid_left_val

        # this is problematic since it doesn't seem that the value at a particular node matters
        assert tree.verify(0, tree.create_merkle_proof(0))

    def test_old(self):
        slot = 2
        txHash = HexBytes('0xcf04ea8bb4ff94066eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84')
        slot2 = 600
        txHash2 = HexBytes('0xabcabcabacbc94566eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84')
        slot3 = 30000
        txHash3 = HexBytes('0xabcaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c9f40d84d5491f5a7735200de010d84')
        tx = {slot: txHash, slot2: txHash2, slot3: txHash3}
        tree = SparseMerkleTree(64, tx)
        for s in tx.keys():
            proof = tree.create_merkle_proof(s)
            inc = tree.verify(s, proof)
            assert inc
