from contract_binds.plasma_cash import PlasmaCash
from config import plasma_config as config
from threading import Thread
import time

abi_file = './abi/RootChain.json'
endpoint = 'http://localhost:8545'

def handle_event(event):
    print(event['args'])

def main():
    bob = "0x4fb2180e2d0c4fdbd7ab22c041aa7faf2e113572"
    plasma = PlasmaCash(config['bob'], abi_file, config['root_chain'], endpoint)
    depositor = { 'depositor' : plasma.account.address } 
    plasma.watch_event('Deposit', handle_event, 1)
    while True:
        time.sleep(1)

main()
