from web3.auto import w3

def sign(hash, key):
    return w3.eth.account.signHash(hash, private_key=key)

def get_sender(hash, sig):
    if sig == None:
        raise InvalidTxSignatureException('Tx not signed')
    return w3.eth.account.recoverHash(hash, signature=sig)
