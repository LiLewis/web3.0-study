// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract Voting {

    //Deploy初始化数据：["Alan", "Bob", "Lewis"]

    // 定义一个 mapping 来存储候选人的得票数
    mapping(string => uint) private  votes;

    // 用于存储已投票的用户
    mapping(address => bool) private hasVoted;

    // 存储候选人列表
    string[] private candidates;

    // 事件：投票成功
    event Voted(address indexed voter, string candidate);

    // 构造函数：初始化候选人列表
    constructor(string[] memory initialCandidates) {
        candidates = initialCandidates;
    }

    // 投票函数
    function vote(string memory candidate) public {
        // 检查用户是否已经投过票
        require(!hasVoted[msg.sender], "You have already voted");

        // 投票给候选人
        votes[candidate] += 1;

        // 标记该用户为已投票
        hasVoted[msg.sender] = true;

        // 触发事件
        emit Voted(msg.sender, candidate);
    }

    // 获取某个候选人的得票数
    function getVotes(string memory candidate) public view returns (uint) {
        return votes[candidate];
    }

    // 重置所有候选人的得票数
    function resetVotes() public {
        // 遍历候选人列表，重置得票数
        for (uint i = 0; i < candidates.length; i++) {
            votes[candidates[i]] = 0;
        }
    }
}