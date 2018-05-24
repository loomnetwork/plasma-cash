from web3.auto import w3
from hexbytes import HexBytes

def sign(hash, key):
    # DO NOT PREFIX!
    sig = HexBytes('0') + w3.eth.account.signHash(hash, private_key=key).signature
    return sig

def get_sender(hash, sig):
    if sig == None:
        raise InvalidTxSignatureException('Tx not signed')
    return w3.eth.account.recoverHash(hash, signature=sig[1:])
