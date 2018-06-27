pragma solidity ^0.4.24;

import "./RootChain.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721.sol";

contract FastWithdrawal {

    RootChain rootChain;

    // Register an exit for a slot, along with the coins
    // that can be used to buyout that exit
    struct Exit {
        address contractAddress;
        address owner;
        bool buyable;
        uint buyoutCoins;
    }
    mapping (uint64 => Exit) exits;

    constructor (RootChain _rootChain) public { rootChain = _rootChain; }

    // We do not add any sanitation checking in this function as
    // it's all done in the plasma contract
    function startExit(
        address contractAddress,
        uint256 buyoutCoins, // to do make an array
        uint64 slot, bytes prevTxBytes,
        bytes exitingTxBytes, bytes prevTxInclusionProof,
        bytes exitingTxInclusionProof, bytes sig,
        uint256[2] blk)
        external
        payable 
    {
        // Start exit on plasma contract
        startExit(slot, prevTxBytes,
            exitingTxBytes, prevTxInclusionProof,
            exitingTxInclusionProof, sig,
            blk
        );

        // Register the exit buyout conditions
        registerExit(slot, contractAddress, buyoutCoins);
    }

    function registerExit(uint64 slot, address contractAddress, uint buyoutCoins) private {
        Exit storage exit = exits[slot];
        exit.owner = msg.sender;
        exit.contractAddress = contractAddress;
        exit.buyoutCoins = buyoutCoins;
        exit.buyable = true;
    }

    // sender must have approved it first!
    // function buyExit


    // hack to avoid stack too deep
    function startExit(
        uint64 slot, bytes prevTxBytes,
        bytes exitingTxBytes, bytes prevTxInclusionProof,
        bytes exitingTxInclusionProof, bytes sig,
        uint256[2] blk)
        private
    {
        rootChain.startExit.value(msg.value)(
            slot, prevTxBytes,
            exitingTxBytes, prevTxInclusionProof,
            exitingTxInclusionProof, sig,
            blk
        );
    }

    function buyExit(
        address contractAddress, 
        uint64 slot,
        uint256 buyoutCoin)
        external
    {
        Exit memory exit = exits[slot];
        exits[slot].owner = msg.sender;
        require(buyoutCoin == exit.buyoutCoins);
        require(contractAddress == exit.contractAddress);
        ERC721(contractAddress).safeTransferFrom(msg.sender, exit.owner, buyoutCoin);
    }

    function withdraw(address contractAddress, uint64 slot, address to) external {
        require(exits[slot].owner == msg.sender);
        rootChain.withdraw(slot, to);
        // ERC721(contractAddress).transferFrom(address(this), msg.sender, uint256(uid));
    }

    function withdrawBonds(address to) external {
        // TODO;
    }

}

