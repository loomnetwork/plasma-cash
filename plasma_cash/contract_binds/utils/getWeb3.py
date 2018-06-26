import time


def getWeb3(endpoint=None):
    if endpoint is None:
        from web3.auto import w3
    elif 'http' in endpoint:
        from web3 import Web3, HTTPProvider

        w3 = Web3(HTTPProvider(endpoint))
    else:
        from web3 import Web3, IPCProvider
        from web3.middleware import geth_poa_middleware

        w3 = Web3(IPCProvider(endpoint))
        w3.middleware_stack.inject(geth_poa_middleware, layer=0)
    # w3.eth.defaultAccount = w3.eth.accounts[0]
    return w3


def waitForTransactionReceipt(w3, tx_hash):
    while True:
        tx_receipt = w3.eth.getTransactionReceipt(tx_hash)
        if tx_receipt is not None:
            break
        time.sleep(1)
    return tx_receipt
