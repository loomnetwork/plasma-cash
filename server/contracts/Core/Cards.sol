// Copyright Loom Network 2018 - All rights reserved, Dual licensed on GPLV3
// Learn more about Loom DappChains at https://loomx.io
// All derivitive works of this code must incluse this copyright header on every file 

pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Token.sol";


contract CryptoCards is ERC721Token("CryptoCards", "CCC") {

    address plasma;

    constructor (address _plasma) public {
        plasma = _plasma;
    }

    function register() external {
        // Give each new player 5 cards
        for (int j = 0; j < 5; j++) {
            create();
        }
    }

    function depositToPlasmaWithData(uint tokenId, bytes _data) public {
        require(plasma != address(0));
        safeTransferFrom(
            msg.sender,
            plasma,
            tokenId,
            _data);
    }

    function depositToPlasma(uint tokenId) public {
        require(plasma != address(0));
        safeTransferFrom(msg.sender, plasma, tokenId);
    }

    function create() private {
        uint256 tokenId = allTokens.length + 1;
        _mint(msg.sender, tokenId);
    }

}
