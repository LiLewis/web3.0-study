// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.10;

contract BeggingContract {
    address public owner;  // 合约所有者

    // 记录每个捐赠者的捐赠金额
    mapping(address => uint256) private donations;

    // 记录所有捐赠者地址（用于排行榜）
    address[] private donors;

    // 捐赠事件
    event Donation(address indexed donor, uint256 amount);

    // 时间限制
    uint256 public donationStart;
    uint256 public donationEnd;

    // 构造函数，设置合约所有者和可选时间限制
    constructor(uint256 _donationStart, uint256 _donationEnd) {
        require(_donationEnd > _donationStart, "End must be after start");
        owner = msg.sender;
        donationStart = _donationStart;
        donationEnd = _donationEnd;
    }

    // 修饰符：限制只有合约所有者可以调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this");
        _;
    }

    // 修饰符：限制在时间段内捐赠（可选）
    modifier withinTime() {
        require(block.timestamp >= donationStart && block.timestamp <= donationEnd, "Donation not allowed now");
        _;
    }

    // 捐赠函数，允许用户发送 ETH
    function donate() public payable withinTime {
        require(msg.value > 0, "Donation must be greater than 0");

        if (donations[msg.sender] == 0) {
            donors.push(msg.sender); // 记录新的捐赠者
        }

        donations[msg.sender] += msg.value;

        emit Donation(msg.sender, msg.value);
    }

    // 查询某个地址的捐赠金额
    function getDonation(address donor) public view returns (uint256) {
        return donations[donor];
    }

    // 合约所有者提取所有捐赠资金
    function withdraw() public onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");

        // 使用call更安全
        (bool success, ) = owner.call{value: balance}("");
        require(success, "Withdraw failed");
    }

    // 获取所有捐赠者地址
    function getAllDonors() public view returns (address[] memory) {
        return donors;
    }

    // 获取排行榜（前 N 个捐赠者）
    function topDonors(uint256 topN) public view returns (address[] memory, uint256[] memory) {
        require(topN > 0, "topN must be greater than 0");

        uint256 n = donors.length < topN ? donors.length : topN;
        address[] memory topAddresses = new address[](n);
        uint256[] memory topAmounts = new uint256[](n);

        // 简单选择排序找出前 N
        address[] memory tempDonors = donors;
        uint256[] memory tempAmounts = new uint256[](tempDonors.length);
        for (uint256 i = 0; i < tempDonors.length; i++) {
            tempAmounts[i] = donations[tempDonors[i]];
        }

        for (uint256 i = 0; i < n; i++) {
            uint256 maxIndex = i;
            for (uint256 j = i + 1; j < tempDonors.length; j++) {
                if (tempAmounts[j] > tempAmounts[maxIndex]) {
                    maxIndex = j;
                }
            }
            // 交换
            (tempDonors[i], tempDonors[maxIndex]) = (tempDonors[maxIndex], tempDonors[i]);
            (tempAmounts[i], tempAmounts[maxIndex]) = (tempAmounts[maxIndex], tempAmounts[i]);

            topAddresses[i] = tempDonors[i];
            topAmounts[i] = tempAmounts[i];
        }

        return (topAddresses, topAmounts);
    }

    // 查看合约余额
    function contractBalance() public view returns (uint256) {
        return address(this).balance;
    }
}
