const HDWalletProvider = require('truffle-hdwallet-provider')
const PrivateKeyProvider = require('truffle-privatekey-provider')
require('babel-register');
require('babel-polyfill');

const mnemonic = 'stumble story behind hurt patient ball whisper art swift tongue ice alien';


let ropstenProvider, kovanProvider, rinkebyProvider = {}

if (process.env.LIVE_NETWORKS) {
  ropstenProvider = new HDWalletProvider(mnemonic, 'https://ropsten.infura.io/')
  kovanProvider = new HDWalletProvider(mnemonic, 'https://kovan.infura.io')

  try {
    const { rpc, key } = require(require('homedir')()+'/.rinkebykey.json')
    rinkebyProvider = new PrivateKeyProvider(key, rpc)
  } catch (e) {
    rinkebyProvider = new HDWalletProvider(mnemonic, 'https://rinkeby.infura.io')
  }
}

const mochaGasSettings = {
  reporter: 'eth-gas-reporter',
  reporterOptions : {
    currency: 'USD',
    gasPrice: 3
  }
}

const mocha = process.env.GAS_REPORTER ? mochaGasSettings : {}

module.exports = {
  networks: {
    rpc: {
      network_id: '*',
      host: 'localhost',
      port: 8545,
    },
    ropsten: {
      network_id: 3,
      provider: ropstenProvider,
      gas: 4700036,
    },
    kovan: {
      network_id: 42,
      provider: kovanProvider,
      gas: 6.9e6,
    },
    rinkeby: {
      network_id: 4,
      provider: rinkebyProvider,
      gas: 6.9e6,
      gasPrice: 15000000001
    },
    coverage: {
      host: "localhost",
      network_id: "*",
      port: 8555,
      gas: 0xffffffffff,
      gasPrice: 0x01
    },
  },
  build: {},
  mocha,
  solc: {
      optimizer: {
          enabled: true,
          runs: 200
      }
  },
}
