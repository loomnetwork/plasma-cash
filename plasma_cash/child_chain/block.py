import rlp
from rlp.sedes import CountableList, binary
from web3.auto import w3

from child_chain.exceptions import (CoinAlreadyIncludedException,
                                    InvalidBlockSignatureException)
from utils.merkle.sparse_merkle_tree import SparseMerkleTree
from utils.utils import get_sender, sign

from .transaction import Transaction


class Block(rlp.Serializable):

    fields = [('transaction_set', CountableList(Transaction)), ('sig', binary)]

    def __init__(self, transaction_set=None, sig=b'\x00' * 65):
        if transaction_set is None:
            self.transactions = {}
        else:
            self.transactions = {tx.uid: tx for tx in transaction_set}
        self.merkle = None
        self.sig = sig

    @property
    def hash(self):
        return w3.sha3(rlp.encode(self, UnsignedBlock))

    @property
    def merkle_hash(self):
        return w3.sha3(rlp.encode(self))

    @property
    def transaction_set(self):
        return list(self.transactions.values())

    @property
    def sender(self):
        if self.sig == b'\x00' * 65:
            raise InvalidBlockSignatureException('Block not signed')
        return get_sender(self.hash, self.sig)

    def merklize_transaction_set(self):
        hashed_transaction_dict = {
            tx.uid: tx.hash for tx in self.transactions.values()
        }
        self.merkle = SparseMerkleTree(64, hashed_transaction_dict)
        return self.merkle.root

    def add_tx(self, tx):
        if tx.uid in self.transactions:
            raise CoinAlreadyIncludedException('double spend rejected')
        else:
            self.transactions[tx.uid] = tx

    # `uid` is the coin that was transferred
    def get_tx_by_uid(self, uid):
        if uid in self.transactions:
            return self.transactions[uid]
        else:
            return Transaction(0, 0, 0, 0)

    def sign(self, key):
        self.sig = sign(self.hash, key)


UnsignedBlock = Block.exclude(['sig'])
