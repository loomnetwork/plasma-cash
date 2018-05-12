const Web3 = require('web3')
const args = require('optimist').string("address").argv;
const truffleContract = require("truffle-contract");

let rpc = args.rpc ? args.arpc : "http://localhost:8545"
let conf = args.abi ? require(args.abi) : require('./build/contracts/RootChain.json');
let abi = conf.abi;
let address = args.address;

const web3 = new Web3(new Web3.providers.HttpProvider(rpc));

const contract = truffleContract({abi});
contract.setProvider(web3.currentProvider);
plasma = contract.at(address);

plasma.Deposit({'depositor': web3.eth.accounts[2]}, {fromBlock: 0, toBlock: 'latest'}).get((error, res) => {
    res.forEach(r => console.log(r.args))
});


