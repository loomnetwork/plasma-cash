const CryptoCards = artifacts.require("CryptoCards");
const RootChain = artifacts.require("RootChain");

module.exports = async function(deployer, network, accounts) {
    // return; // for testing
    let aCryptoCardsInstance;
    let aRootChainInstance;

    return deployer.deploy(RootChain)
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

            aRootChainInstance.setCryptoCards(instance.address);
        });
};

