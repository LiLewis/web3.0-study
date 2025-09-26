// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

// import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

//用于模拟 ERC20 代币在拍卖中作为支付货币。
//可以在 bidERC20() 中测试 ERC20 出价功能。

//声明一个 ERC20 合约，同时支持 permit(导入 ERC20Permit 扩展，实现 EIP-2612 permit（离链签名授权）)
// 模拟 ERC20 代币，在拍卖中作为支付货币
contract TestERC20 is ERC20Permit {
    //创建一个构造函数
    //设置 token 名称、符号，并 mint 初始代币
    constructor() ERC20("MyToken", "MTK") ERC20Permit("MyToken") {
        // 部署者获得 100,000 MTK（18 位精度）
        _mint(msg.sender, 100_000 * 10 ** decimals());
    }
}