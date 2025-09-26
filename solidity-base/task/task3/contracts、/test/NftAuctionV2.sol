// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "./NftAuction.sol";

//拓展 NftAuction，用于测试 UUPS 升级是否生效。
//可以在升级后调用 testHello() 验证代理升级成功。
contract NftAuctionV2 is NftAuction {
    // 新增示例方法用于验证升级生效
    function testHello() external pure returns (string memory) {
        return "hello from V2!";
    }
}