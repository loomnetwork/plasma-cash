# Plasma Cash ERC721

Root Chain contract implements ERC721 receiver and allows receiving ONLY erc721 from a contract that it gets connected with. 

When it receives a transfer from the connected ERC721 contract, it calls `deposit` and emits an event to be watched by the client.
