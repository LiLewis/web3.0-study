// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract BinarySearch {

    //Deploy 输入：arr = [1, 2, 3, 4, 5, 6]    target = 1~6 or 7

    // 二分查找函数，查找目标值在有序数组中的索引
    function binarySearch(uint[] memory arr, uint target) public pure returns (int) {
        uint left = 0;
        uint right = arr.length - 1;

        // 进行二分查找
        while (left <= right) {
            uint mid = left + (right - left) / 2;

            // 检查中间元素
            if (arr[mid] == target) {
                // 返回目标值的索引
                return int(mid);
            }

            // 如果目标值小于中间元素，缩小查找范围到左半部分
            if (arr[mid] > target) {
                right = mid - 1;
            }
                // 如果目标值大于中间元素，缩小查找范围到右半部分
            else {
                left = mid + 1;
            }
        }

        // 如果没有找到目标值，返回 -1
        return -1;
    }
}