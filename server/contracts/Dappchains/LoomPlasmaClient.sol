pragma solidity ^0.4.21;

import "./Dappchain.sol";

//This file is for solidity contracts running on Loom DappChains 
//It should not be run on a vanilla Ethereum instance

contract LoomPlasma is DappChain {

    function senduxto(bytes _input) public  {
        loomtransferUxto(_input); 
    }

}