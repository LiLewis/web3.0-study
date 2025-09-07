// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721URIStorage, Ownable {
    // NFT 自增 ID
    uint256 private _tokenIds;

    // 构造函数必须给 ERC721 和 Ownable 传参
    constructor() ERC721("MyNFT", "MNFT") Ownable(msg.sender) {}

    /// @notice 铸造 NFT
    /// @param recipient 接收 NFT 的钱包地址
    /// @param tokenURI 存储在 IPFS 的元数据 JSON 链接
    function mintNFT(address recipient, string memory tokenURI)
    public
    onlyOwner
    returns (uint256)
    {
        _tokenIds++;
        uint256 newItemId = _tokenIds;

        _mint(recipient, newItemId);
        _setTokenURI(newItemId, tokenURI);

        return newItemId;
    }
}

//metadata.json代码
//{
//    "name": "My First NFT",
//    "description": "这是我在 Sepolia 测试网上铸造的第一个 NFT",
//    "image": "ipfs://xxxxxxx.png",
//    "attributes": [
//        {
//        "trait_type": "Coolness",
//        "value": 100
//        }
//    ]
//}
