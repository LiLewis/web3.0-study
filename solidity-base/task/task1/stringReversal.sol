// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract stringReversal {

    //Deploy参数例子：abcdef

    // 反转字符串的函数
    function reverseString(string memory input) public pure returns (string memory) {
        // 将字符串转换为字节数组
        bytes memory str = bytes(input);
        // 获取字符串长度
        uint len = str.length;
        // 创建一个新的字节数组来存储反转后的字符串
        bytes memory reversed = new bytes(len);

        // 反转字符串
        for (uint i = 0; i < len; i++) {
            // 从后向前填充字节
            reversed[i] = str[len - 1 - i];
        }
        // 返回result
        return string(reversed);
    }
}