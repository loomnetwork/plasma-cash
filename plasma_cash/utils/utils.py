from hexbytes import HexBytes
from web3.auto import w3

from child_chain.exceptions import InvalidTxSignatureException


def sign(hash, key):
    # DO NOT PREFIX!
    sig = (
        HexBytes('0')
        + w3.eth.account.signHash(hash, private_key=key).signature
    )
    return sig


def get_sender(hash, sig):
    if sig is None:
        raise InvalidTxSignatureException('Tx not signed')
    return w3.eth.account.recoverHash(hash, signature=sig[1:])


def increaseTime(w3, time):
    start = w3.eth.getBlock('latest').timestamp
    # provider.make_request(method='evm_increaseTime', params=start+time)
    w3.manager.request_blocking(
        method='evm_increaseTime', params=[start + time]
    )
