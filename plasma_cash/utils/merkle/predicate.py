from ethereum.utils import sha3


def is_valid_proof(leaf, uid, proof, root):
    index = uid
    computed_hash = leaf
    for idx in range(0, len(proof), 32):
        segment = proof[idx:idx + 32]
        if index % 2 == 0:
            computed_hash = sha3(computed_hash + segment)
        else:
            computed_hash = sha3(segment + computed_hash)
        index = index // 2
    return computed_hash == root
