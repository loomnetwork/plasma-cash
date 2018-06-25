const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");
const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");

module.exports = async function(deployer, network, accounts) {
    // return; // for testing
    let aCryptoCardsInstance;
    let aRootChainInstance;
    let aValidatorManagerContractInstance;

    return deployer.deploy(ValidatorManagerContract)
        .then(() => ValidatorManagerContract.deployed())
        .then(instance => {
            aValidatorManagerContractInstance = instance;
            console.log('ValidatorManagerContract deployed at address: ' + instance.address);
            return deployer.deploy(RootChain, instance.address);
        })
        .then(() => RootChain.deployed())
        .then(instance => {
            aRootChainInstance = instance;
            console.log('RootChain deployed at address: ' + instance.address);
            return deployer.deploy(CryptoCards, instance.address);
        })
    .then(() => CryptoCards.deployed())
        .then((instance) => {
            aCryptoCardsInstance = instance;
            console.log('CryptoCards deployed at address: ' + instance.address);

            aValidatorManagerContractInstance.toggleToken(instance.address);
        });
};

