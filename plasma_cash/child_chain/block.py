import rlp
from rlp.sedes import CountableList, binary
from web3.auto import w3

from .transaction import Transaction
from utils.utils import get_sender, sign
from utils.merkle.sparse_merkle_tree import SparseMerkleTree

class Block(rlp.Serializable):

    fields = [
        ('transaction_set', CountableList(Transaction)),
        ('sig', binary)
    ]

    def __init__(self, transaction_set=None, sig=b'\x00' * 65):
        if transaction_set is None:
            transaction_set = []
        self.transaction_set = transaction_set
        self.merkle = None
        self.sig = sig

    @property
    def hash(self):
        return w3.sha3(rlp.encode(self, UnsignedBlock))

    @property
    def merkle_hash(self):
        return w3.sha3(rlp.encode(self))

    @property
    def sender(self):
        if self.sig == b'\x00' * 65:
            raise InvalidBlockSignatureException('Block not signed')
        return get_sender(self.hash, self.sig)

    def merklize_transaction_set(self):
        hashed_transaction_dict = {tx.uid: tx.hash for tx in self.transaction_set}
        self.merkle = SparseMerkleTree(64, hashed_transaction_dict)
        return self.merkle.root

    def add_tx(self, tx):
        self.transaction_set.append(tx)

    # `uid` is the coin that was transferred
    def get_tx_by_uid(self, uid):
        for tx in self.transaction_set: # replace with better searching 
            if tx.uid == uid:
                return tx
        return None


    def sign(self, key):
        self.sig = sign(self.hash, key).signature

UnsignedBlock = Block.exclude(['sig'])
