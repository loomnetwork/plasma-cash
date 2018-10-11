const HDWalletProvider = require('truffle-hdwallet-provider')
const PrivateKeyProvider = require('truffle-privatekey-provider')
require('babel-register');
require('babel-polyfill');

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
      provider: () => new HDWalletProvider(process.env.mnemonic, 'https://ropsten.infura.io'),
      gas: 4700036,
    },
    kovan: {
      network_id: 42,
      provider: () => new HDWalletProvider(process.env.mnemonic, 'https://kovan.infura.io'),
      gas: 6.9e6,
    },
    rinkeby: {
      network_id: 4,
      provider: () => new HDWalletProvider(process.env.mnemonic, 'https://rinkeby.infura.io'),
      gas: 6.9e6,
      skipDryRun: true,
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
