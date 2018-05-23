from sparse_merkle_tree import *
from hexbytes import HexBytes
slot = 2
txHash = HexBytes('0xcf04ea8bb4ff94066eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84')
slot2 = 600
txHash2 = HexBytes('0xabcabcabacbc94566eb84dd932f9e66d1c9f40d84d5491f5a7735200de010d84')

slot3 = 30000
txHash3 = HexBytes('0xabcaaaaaaaaaaaaaaaaaaaaaaaaaaaaa1c9f40d84d5491f5a7735200de010d84')

tx = {slot:txHash, slot2:txHash2, slot3:txHash3}
tree = SparseMerkleTree(64, tx);
for s in tx.keys():
    proof = tree.create_merkle_proof(s)
    inc = tree.verify(s, proof)
    assert inc == True
print ('OK! SMT WORKS :) ')
