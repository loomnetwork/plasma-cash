from web3.auto import w3


# TODO this is not necessary. I'll remove this ASAP and delete this file once
# we're using the generic sha3 fuction in sparse_merkle_tree.py
def keccak256(*args):
    hashes = list(args)
    bytes32 = ['bytes32'] * len(hashes)
    return w3.soliditySha3(bytes32, hashes)
