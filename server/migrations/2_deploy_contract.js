const CryptoCards = artifacts.require("CryptoCards");
const LoomToken = artifacts.require("LoomToken")
const RootChain = artifacts.require("RootChain");
const ValidatorManagerContract = artifacts.require("ValidatorManagerContract");

module.exports = async function(deployer, network, accounts) {

   //  deployer.deploy(ValidatorManagerContract).then(async () => {
   //      const vmc = await ValidatorManagerContract.deployed();
   //      console.log(`ValidatorManagerContract deployed at address: ${vmc.address}`);

   await deployer.deploy(RootChain, "0x759cEf0AE8855f8fAd5e74d91c3590639a6451eC");
   const root = await RootChain.deployed();
   console.log(`RootChain deployed at address: ${root.address}`);

   // await deployer.deploy(LoomToken, root.address);
   // const erc20 = await LoomToken.deployed();
   // console.log(`LoomToken deployed at address: ${erc20.address}`);

    // const decimals = 10 **18
    // await erc20.transfer("0x3d5cf1f50c7124acbc6ea69b96a912fe890619d0", 500 * decimals, {from: "0xC5dFc9282BF68DFAd041a04a0c09bE927b093992"})

   //      await deployer.deploy(CryptoCards, root.address);
   //      const cards = await CryptoCards.deployed();
   //      console.log(`CryptoCards deployed at address: ${cards.address}`);

   //      await vmc.toggleToken(cards.address);
   //  });
    //     await deployer.deploy(LoomToken, root.address);
    //     const erc20 = await CryptoCards.deployed();
    //     console.log(`Loom deployed at address: ${erc20.address}`);

    //     await vmc.toggleToken(cards.address);
    // });
};

