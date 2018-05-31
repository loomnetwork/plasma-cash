# Plasma Cash ERC721

`Cards.sol` represents an ERC721 implementation of an NFT contract. 

`RootChain.sol` represents the Plasma Cash Contract. The contract inherits from `ERC721Receiver` in order to trigger the `deposit` function, when an NFT gets transferred to the address of the contract. It also inherits from `SparseMerkleTree` in order to do validation of Merkle Proofs.

# Features
- Blocks that are not multiple of `interval=1000` have smaller proofs
- Merkle Proof Validation for other blocks
- Double Spend Challenge
- Exit Spent Coin Challenge
- Exit with Invalid history Challenge & Response

