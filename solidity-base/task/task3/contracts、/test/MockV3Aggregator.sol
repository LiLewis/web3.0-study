// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8;

import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

//拍卖合约 _toUsdScaled 函数用来计算出价美元价值。
//可以测试价格变化、ERC20/ETH 出价逻辑。

//模拟 Chainlink Price Feed，用于本地测试。
contract MockV3Aggregator is AggregatorV3Interface {
    //位数
    uint8 private _decimals;
    //存储最新价格（如 ETH/USD）
    int256 private _latestAnswer;

    constructor(uint8 decimals_, int256 initialAnswer_) {
        _decimals = decimals_;
        _latestAnswer = initialAnswer_;
    }

    function decimals() external view override returns (uint8) {
        return _decimals;
    }

    function description() external pure override returns (string memory) {
        return "MockV3Aggregator";
    }

    function version() external pure override returns (uint256) {
        return 0;
    }

    //返回价格数据
    function getRoundData(uint80)
        external
        view
        override
        returns (
            uint80,
            int256,
            uint256,
            uint256,
            uint80
        )
    {
        return (0, _latestAnswer, block.timestamp, block.timestamp, 0);
    }

    //返回价格数据
    function latestRoundData()
        external
        view
        override
        returns (
            uint80,
            int256,
            uint256,
            uint256,
            uint80
        )
    {
        return (0, _latestAnswer, block.timestamp, block.timestamp, 0);
    }

    //可手动修改价格
    function updateAnswer(int256 newAnswer) external {
        _latestAnswer = newAnswer;
    }
}