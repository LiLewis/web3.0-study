// SPDX-License-Identifier: MIT
pragma solidity ^0.8;

contract MergeArray {

    //Deploy 输入：nums1 = [1, 3, 5]    nums2 = [2, 4, 6]

    // 合并两个有序数组，并返回新的有序数组
    function merge(uint[] memory nums1, uint[] memory nums2) public pure returns (uint[] memory) {
        uint m = nums1.length;
        uint n = nums2.length;

        // 创建一个新数组存储合并后的结果
        uint[] memory merged = new uint[](m + n);

        uint i = 0;
        uint j = 0;
        uint k = 0;

        // 通过双指针方法合并两个有序数组
        while (i < m && j < n) {
            if (nums1[i] <= nums2[j]) {
                merged[k++] = nums1[i++];
            } else {
                merged[k++] = nums2[j++];
            }
        }

        // 如果 nums1 中还有剩余，直接放到结果数组中
        while (i < m) {
            merged[k++] = nums1[i++];
        }

        // 如果 nums2 中还有剩余，直接放到结果数组中
        while (j < n) {
            merged[k++] = nums2[j++];
        }

        return merged;
    }
}