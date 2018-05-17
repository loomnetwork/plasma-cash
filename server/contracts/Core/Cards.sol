pragma solidity ^0.4.22;

import 'openzeppelin-solidity/contracts/token/ERC721/ERC721Token.sol';

contract CryptoCards is ERC721Token("CryptoCards", "CCC") {

    mapping(address => bool) private registered;

    address plasma;

    constructor (address _plasma) public {
        plasma = _plasma;
    }

    function register() external {
        // require(!registered[msg.sender]);
        for (int j = 0; j<5 ; j++) {
            create(); // Give each new player 5 cards
        }
        // registered[msg.sender] = true;
    }

	function create() private {
		uint256 tokenId = allTokens.length + 1;
		_mint(msg.sender, tokenId);
	}

    function depositToPlasmaWithData(uint tokenId, bytes _data) public {
        require(plasma != address(0));
        safeTransferFrom(msg.sender, plasma, tokenId, _data);
    }

    function depositToPlasma(uint tokenId) public {
        require(plasma != address(0));
        safeTransferFrom(msg.sender, plasma, tokenId);
    }

}
