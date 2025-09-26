// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

//用于生成测试 NFT，在拍卖中上架。
//可以测试 NFT 转移、拍卖创建、结束等逻辑。

//ERC721 扩展，可给 NFT 单独设置 URI,    提供 onlyOwner 修饰器控制权限
contract TestERC721 is ERC721URIStorage, Ownable {
    // 初始化：传入合约 owner 地址
    constructor() ERC721("Troll", "TROLL") Ownable(msg.sender) {}

    // 铸造 NFT，只允许 owner 调用
    function mint(address to, uint256 tokenId) external onlyOwner {
        _safeMint(to, tokenId);
    }

    // 设置 NFT 元数据 URI
    function setTokenURI(uint256 tokenId, string memory newTokenURI) external onlyOwner {
        _setTokenURI(tokenId, newTokenURI);
    }
}