import rlp
import ethereum.utils
from web3.auto import w3
from rlp.sedes import big_endian_int, binary

from utils.utils import get_sender, sign
from .exceptions import InvalidTxSignatureException

class Transaction(rlp.Serializable):

    fields = [
        ('prev_block', big_endian_int),
        ('uid', big_endian_int),
        ('new_owner', ethereum.utils.address),
        ('sig', binary)
    ]

    def __init__(self, prev_block, uid, new_owner, sig=b'\x00' * 65):
        self.prev_block = prev_block
        self.uid = uid
        self.new_owner = ethereum.utils.normalize_address(new_owner)
        self.sig = sig
        self.spent = False # not part of the rlp
        self.make_mutable()

    @property
    def hash(self):
        return w3.sha3(rlp.encode(self, UnsignedTransaction))

    @property
    def merkle_hash(self):
        return w3.sha3(self.hash + self.sig)

    @property
    def sender(self):
        if self.sig == b'\x00' * 65:
            raise InvalidTxSignatureException('Tx not signed')
        return get_sender(self.hash, self.sig)

    def sign(self, key):
        self.sig = sign(self.hash, key).signature

UnsignedTransaction = Transaction.exclude(['sig'])
