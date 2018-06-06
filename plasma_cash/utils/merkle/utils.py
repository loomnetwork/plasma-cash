from web3.auto import w3


def keccak256(*args):
    hashes = list(args)
    bytes32 = ['bytes32'] * len(hashes)
    return w3.soliditySha3(bytes32, hashes)
