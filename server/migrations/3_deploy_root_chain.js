const CryptoCards = artifacts.require("CryptoCards");
const LoomToken = artifacts.require("LoomToken")
const RootChain = artifacts.require("RootChain");

module.exports = async function(deployer, network, accounts) {
    let vmcAddress = '0xf1ffd3f1598cd054f5031b2af219d85b0d443175';
    console.log(`ValidatorManagerContract deployed at address: ${vmcAddress}`);

    await deployer.deploy(RootChain, vmcAddress);
    const root = await RootChain.deployed();
    console.log(`RootChain deployed at address: ${root.address}`);        
};



