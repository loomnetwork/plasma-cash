import binascii
import pytest
from ethereum.utils import sha3
from utils.merkle.sparse_merkle_tree import SparseMerkleTree

class TestSparseMerkleTree(object):
    def test_emptySMT(self):
        emptyTree = SparseMerkleTree(64, {})
        assert len(emptyTree.leaves) == 0

    def test_twoItemSMT(self):
        twoItemSMT = SparseMerkleTree(64, { 1: bytearray.fromhex('46fc5f183973cb120cfde52cc65295b4429034c3036f9a16260d02454c707b5c')
                                          , 2: bytearray.fromhex('598219f259df021d8251187348dcd17c46ff86a42286b43d4660c1e19d160c2d')
                                          })
        assert len(twoItemSMT.leaves) == 2
        assert twoItemSMT.create_merkle_proof(1) == bytearray.fromhex('00000000000000028711aeb008ae8a08665f7d8cbeac73dee60e2d82543b11df0f05406eee3baef8')

#   def test_all_leaves_with_val(self):
#       dummy_val = b'\x01' * 32
#       leaves = {0: dummy_val, 1: dummy_val, 2: dummy_val, 3: dummy_val}
#       tree = SparseMerkleTree(depth=3, leaves=leaves)
#       mid_level_val = sha3(dummy_val + dummy_val)
#       assert tree.root == sha3(mid_level_val + mid_level_val)
#
#   def test_empty_leaves(self):
#       tree = SparseMerkleTree(depth=3)
#       empty_val = b'\x00' * 32
#       mid_level_val = sha3(empty_val + empty_val)
#       assert tree.root == sha3(mid_level_val + mid_level_val)

#   def test_empty_left_leave(self):
#       empty_val = b'\x00' * 32
#       dummy_val = b'\x01' * 32
#       leaves = {1: dummy_val, 2: dummy_val, 3: dummy_val}
#       tree = SparseMerkleTree(depth=3, leaves=leaves)
#       mid_left_val = sha3(empty_val + dummy_val)
#       mid_right_val = sha3(dummy_val + dummy_val)
#       assert tree.root == sha3(mid_left_val + mid_right_val)
#
#   def test_empty_right_leave(self):
#       empty_val = b'\x00' * 32
#       dummy_val = b'\x01' * 32
#       leaves = {0: dummy_val, 2: dummy_val, 3: dummy_val}
#       tree = SparseMerkleTree(depth=3, leaves=leaves)
#       mid_left_val = sha3(dummy_val + empty_val)
#       mid_right_val = sha3(dummy_val + dummy_val)
#       assert tree.root == sha3(mid_left_val + mid_right_val)
#
    def test_exceed_tree_size(self):
        with pytest.raises(SparseMerkleTree.TreeSizeExceededException):
            SparseMerkleTree(depth=1, leaves={0: '0', 1: '1'})

#   def test_create_merkle_proof(self):
#       empty_val = b'\x00' * 32
#       dummy_val = b'\x01' * 32
#       leaves = {0: dummy_val, 2: dummy_val, 3: dummy_val}
#       tree = SparseMerkleTree(depth=3, leaves=leaves)
#       mid_left_val = sha3(dummy_val + empty_val)
#       mid_right_val = sha3(dummy_val + dummy_val)
#       assert tree.create_merkle_proof(0) == empty_val + mid_right_val
#       assert tree.create_merkle_proof(1) == dummy_val + mid_right_val
#       assert tree.create_merkle_proof(2) == dummy_val + mid_left_val
#       assert tree.create_merkle_proof(3) == dummy_val + mid_left_val


# TODO ressurect this old test
# slot = 2
# txHash = HexBytes('0xcf04ea8bb4ff94066eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84')
# slot2 = 600
# txHash2 = HexBytes('0xabcabcabacbc94566eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84')
#
# slot3 = 30000
# txHash3 = HexBytes('0xabcaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c9f40d84d5491f5a7735200de010d84')
#
# tx = {slot: txHash, slot2: txHash2, slot3: txHash3}
# tree = SparseMerkleTree(64, tx)
# for s in tx.keys():
#     proof = tree.create_merkle_proof(s)
#     inc = tree.verify(s, proof)
#     assert inc
# print('OK! SMT WORKS :) ')
