// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract RomanConverter {
    // 整数到罗马数字的转换函数
    function intToRoman(uint num) public pure returns (string memory) {
        // 定义罗马数字符号与整数值的对应关系
        uint16[13] memory values = [1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1];
        string[13] memory symbols = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];

        string memory result = "";

        // 逐步减去数值，并将符号拼接到结果中
        for (uint i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                num -= values[i];
                result = string(abi.encodePacked(result, symbols[i]));
            }
        }

        return result;
    }

    // 罗马数字到整数的转换函数
    function romanToInt(string memory s) public pure returns (int) {
        int result = 0;
        uint len = bytes(s).length;

        // 遍历每个字符
        for (uint i = 0; i < len; i++) {
            bytes1 currentChar = bytes(s)[i];
            int currentValue = getRomanValue(currentChar);

            // 获取下一个字符的值
            // int nextValue = (i + 1 < len) ? getRomanValue(bytes(s)[i + 1]) : int(0); // 显式转换0为int类型

            // 如果当前字符小于下一个字符，则减去当前字符的值
            // if (currentValue < nextValue) {
            //     result -= currentValue;
            // } else {
            //     result += currentValue;
            // }

            result += currentValue;
        }


        return result;
    }

    // 获取罗马字符的对应整数值
    function getRomanValue(bytes1 romanChar) private pure returns (int) {
        if (romanChar == 'I') return 1;
        if (romanChar == 'V') return 5;
        if (romanChar == 'X') return 10;
        if (romanChar == 'L') return 50;
        if (romanChar == 'C') return 100;
        if (romanChar == 'D') return 500;
        if (romanChar == 'M') return 1000;
        return 0; // 无效字符
    }
}