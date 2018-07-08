import pytest
from eth_utils.crypto import keccak
from hexbytes import HexBytes

from utils.merkle.sparse_merkle_tree import SparseMerkleTree

empty_val = b'\x00' * 32
default_hash = keccak(empty_val)
dummy_val = keccak(2)
dummy_val_2 = keccak(3)


class TestSparseMerkleTree(object):
    def test_size_limits(self):
        with pytest.raises(SparseMerkleTree.TreeSizeExceededException):
            SparseMerkleTree(depth=0, leaves={0: '0', 1: '1'})
        with pytest.raises(SparseMerkleTree.TreeSizeExceededException):
            SparseMerkleTree(
                depth=1, leaves={0: empty_val, 1: empty_val, 2: empty_val}
            )

    def test_empty_SMT(self):
        emptyTree = SparseMerkleTree(64, {})
        assert len(emptyTree.leaves) == 0
        assert (
            emptyTree.root
            == bytes(HexBytes('0x6f35419d1da1260bc0f33d52e8f6d73fc5d672c0dca13bb960b4ae1adec17937'))
        )

    def test_all_leaves_with_val(self):
        leaves = {0: dummy_val, 1: dummy_val, 2: dummy_val, 3: dummy_val}
        tree = SparseMerkleTree(depth=2, leaves=leaves)
        mid_level_val = keccak(dummy_val * 2)
        assert tree.root == keccak(mid_level_val + mid_level_val)

    def test_empty_leaves(self):
        tree = SparseMerkleTree(depth=2)
        mid_level_val = keccak(default_hash * 2)
        assert tree.root == keccak(mid_level_val * 2)

    def test_empty_left_leave(self):
        leaves = {1: dummy_val, 2: dummy_val, 3: dummy_val}
        tree = SparseMerkleTree(depth=2, leaves=leaves)
        mid_left_val = keccak(default_hash + dummy_val)
        mid_right_val = keccak(dummy_val * 2)
        assert tree.root == keccak(mid_left_val + mid_right_val)

    def test_empty_right_leave(self):
        leaves = {0: dummy_val, 2: dummy_val, 3: dummy_val}
        tree = SparseMerkleTree(depth=2, leaves=leaves)
        mid_left_val = keccak(dummy_val + default_hash)
        mid_right_val = keccak(dummy_val * 2)
        assert tree.root == keccak(mid_left_val + mid_right_val)

    def test_create_merkle_proof(self):
        leaves = {0: dummy_val, 2: dummy_val, 3: dummy_val_2}
        tree = SparseMerkleTree(depth=2, leaves=leaves)
        mid_left_val = keccak(dummy_val + default_hash)
        mid_right_val = keccak(dummy_val + dummy_val_2)
        assert (
            tree.create_merkle_proof(0)
            == (2).to_bytes(8, byteorder='big') + mid_right_val
        )
        assert (
            tree.create_merkle_proof(1)
            == (3).to_bytes(8, byteorder='big') + dummy_val + mid_right_val
        )
        assert (
            tree.create_merkle_proof(2)
            == (3).to_bytes(8, byteorder='big') + dummy_val_2 + mid_left_val
        )
        assert (
            tree.create_merkle_proof(3)
            == (3).to_bytes(8, byteorder='big') + dummy_val + mid_left_val
        )

    def test_old(self):
        slot = 2
        txHash = HexBytes(
            '0xcf04ea8bb4ff94066eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84'
        )
        slot2 = 600
        txHash2 = HexBytes(
            '0xabcabcabacbc94566eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84'
        )
        slot3 = 30000
        txHash3 = HexBytes(
            '0xabcaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c9f40d84d5491f5a7735200de010d84'
        )
        tx = {slot: txHash, slot2: txHash2, slot3: txHash3}
        tree = SparseMerkleTree(64, tx)
        for s in tx.keys():
            proof = tree.create_merkle_proof(s)
            inc = tree.verify(s, proof)
            assert inc

    def test_real_slot_proofs(self):
        slot = 14414645988802088183
        txHash = HexBytes('0x510a183d5457e0d22951440a273f0d8e28e01d15f750d79fd1b27442299f7220')
        tree = SparseMerkleTree(64, {slot: txHash})
        proof = tree.create_merkle_proof(slot)
        inc = tree.verify(slot, proof)
        assert inc

    def test_real_tree_roots(self):
        slot = 14414645988802088183
        txHash = HexBytes('0x4b114962ecf0d681fa416dc1a6f0255d52d701ab53433297e8962065c9d439bd')
        tree = SparseMerkleTree(64, {slot: txHash})
        assert tree.root == bytes(HexBytes('0x0ed6599c03641e5a20d9688f892278dbb48bbcf8b1ff2c9a0e2b7423af831a83'))

        slot = 14414645988802088183
        txHash = HexBytes('0x510a183d5457e0d22951440a273f0d8e28e01d15f750d79fd1b27442299f7220')
        tree = SparseMerkleTree(64, {slot: txHash})
        assert tree.root == bytes(HexBytes('0x8d0ae4c94eaad54df5489e5f9d62eeb4bf06ff774a00b925e8a52776256e910f'))
