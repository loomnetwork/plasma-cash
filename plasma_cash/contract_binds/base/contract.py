from ..utils.getWeb3 import getWeb3
from ..utils.colors import green, red, yellow
from threading import Thread
import time
import json
from web3.utils.events import get_event_data


class Contract(object):
    '''Base class for interfacing with a contract'''
    def __init__(self, keystore, address, abi_file, endpoint):
        w3 = getWeb3(endpoint)
        # address = w3.toChecksumAddress(address) # need to comment out because of ganache bug
        with open(abi_file) as f:
            abi = json.load(f)['abi']
        contract = w3.eth.contract(abi=abi, address=address)
        self.w3 = w3
        # self.w3.eth.enable_unaudited_features()
        if keystore is not None:
            self.account = self.to_account(keystore)
            self.nonce = self.w3.eth.getTransactionCount(self.account.address)
        self.contract = contract

    def to_account(self, data):
        account = self.w3.eth.account.privateKeyToAccount(data)
        del data
        return account

    def sign_and_send(self, func, args, value=0, gas=1000000):
        ''' Expecting all arguments in 1 array '''
        signed_tx = self._sign_function_call(
                func,
                args,
                value,
                gas  # may need to change gas
                )
        self.nonce += 1  # Increment nonce after signing a tx

        info = '{}, Args: {}'.format(func.__name__, args)
        try:
            tx_hash = self._send_raw_tx(signed_tx)
            # info = 'Success '+info
            # self.logger.debug(green(info))
        except Exception as e:
            print(e)
            info = 'Failed '+info
            # self.logger.debug(red(info))
        return tx_hash

    def send_transaction(self, to, value):
        signed_tx = self._sign_transaction(to, value)
        self.nonce += 1  # Increment nonce after signing a tx
        tx_hash = self._send_raw_tx(signed_tx)

        info = 'Sent {} to {}'.format(value, to)
        # self.logger.info(green(info))
        return tx_hash

    def _sign_transaction(self, to, value):
        gas = 21000
        gasPrice = self.w3.toWei('10', 'gwei')

        raw_tx = {
                'chainId': int(self.w3.version.network),
                'to': self.w3.toChecksumAddress(to),
                'value': value,
                'gas': gas,
                'gasPrice': gasPrice,
                # 'nonce': self.nonce
                'nonce': self.w3.eth.getTransactionCount(self.account.address)
                }

        # print(raw_tx)

        signed_tx = self.account.signTransaction(raw_tx)

        return signed_tx

    def _sign_function_call(self, func, args, value, gas):
        """
            Takes reading and timestamp and creates a
            raw transaction call to `ping` at the target contract
            TODO: Add option to modify gas
        """
        info = 'Args => {}'.format(args)
        # self.logger.debug(yellow(info))
        # Build the raw transaction

        raw_tx = func(*args).buildTransaction({
            'gas': gas,
            'value': value,
            # 'gasPrice': self.w3.toWei('10', 'gwei'),
            # 'nonce': self.nonce
            'nonce': self.w3.eth.getTransactionCount(self.account.address)

            })
        raw_tx['to'] = self.w3.toChecksumAddress(raw_tx['to'])

        # info = 'Raw transaction before signing => {}'.format(raw_tx)
        # self.logger.debug(yellow(info))

        # Sign the transaction with the meter's private key
        signed_tx = self.account.signTransaction(raw_tx)

        # info = 'Raw transaction after signing => {}'.format(signed_tx)
        # self.logger.debug(yellow(info))

        return signed_tx

    def _send_raw_tx(self, signed_tx):
        tx_hash = self.w3.eth.sendRawTransaction(signed_tx.rawTransaction)
        return tx_hash
    # Wait until the transaction gets mined
        # receipt = self.waitForTxReceipt(tx)
        # info = 'Transaction Receipt => {}'.format(receipt)
        # self.logger.debug(yellow(info))

    def waitForTxReceipt(self, tx):
        receipt = self.w3.eth.getTransactionReceipt(tx)
        while receipt is None:
            info = 'Waiting for transaction to get mined...'
            # self.logger.debug(yellow(info))
            time.sleep(12)  # Block time avg
            receipt = self.w3.eth.getTransactionReceipt(tx)
        return receipt

    def get_event_data(self, event_name, tx_hash):
        tx_logs = self.w3.eth.getTransactionReceipt(tx_hash)['logs']
        event_abi = self.contract._find_matching_event_abi('Deposit')
        matched = []
        for log in tx_logs:
            try:
                d = get_event_data(event_abi, log)
            except: 
                continue
            matched.append(d)
        return matched

    def watch_event(self, event_name, callback, interval, fromBlock=0,
                toBlock='latest', filters=None):
        event_filter = self.install_filter(
                event_name,
                fromBlock,
                toBlock,
                filters
            )
        worker = Thread(target=self.watcher,
                        args=(event_filter, callback, interval),
                        daemon=True).start()

    def watcher(self, event_filter, callback, interval):
        while True:
            for event in event_filter.get_new_entries():
                callback(event)
                time.sleep(interval)

    def install_filter(self, event_name, fromBlock=0, toBlock='latest',
                       filters=None):
        event = getattr(self.contract.events, event_name)
        eventFilter = event.createFilter(fromBlock=fromBlock,
                                         toBlock=toBlock,
                                         argument_filters=filters)
        return eventFilter
