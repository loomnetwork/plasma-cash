// Copyright Loom Network 2018 - All rights reserved, Dual licensed on GPLV3
// Learn more about Loom DappChains at https://loomx.io
// All derivitive works of this code must incluse this copyright header on every file 

pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/ownership/Ownable.sol";


contract ValidatorManagerContract is Ownable {

    mapping (address => bool) public validators;
    mapping (address => bool) public allowedTokens;

    function checkValidator(address _address) public view returns (bool) {
        // owner is a permanent validator
        if (_address == owner)
            return true;
        return validators[_address];
    }

    function toggleValidator(address _address) public onlyOwner {
        validators[_address] = !validators[_address];
    }

    function toggleToken(address _token) public {
        require(checkValidator(msg.sender), "not a validator");
        allowedTokens[_token] = !allowedTokens[_token];
    }

}
