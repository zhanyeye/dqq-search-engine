package playground

// IsBitOne 判断第i位上是否为1，i从1开始
func IsBitOne(number uint64, bitPosition uint) bool {
	if bitPosition < 1 || bitPosition > 64 {
		panic("index out of range")
	}
	return number&(1<<(bitPosition-1)) != 0
}

// UpdateBitToOne 将第i位设置为1, i 从1开始
func UpdateBitToOne(number uint64, bitPosition uint) uint64 {
	if bitPosition < 1 || bitPosition > 64 {
		panic("index out of range")
	}
	number |= 1 << (bitPosition - 1)
	return number
}

// CountBitOne 统计一个数的二进制有几个1
func CountBitOne(number uint64) int {
	cnt := 0
	c := uint64(1)
	for i := 0; i < 64; i++ {
		if number&c == c {
			cnt++
		}
		c <<= 1
	}
	return cnt
}

/*
CountBitOnePlus 使用 Brian Kernighan 算法计算一个无符号整数中 1 的数量。

算法概要：
1. 每次通过将数字减去 1 与原数字进行按位与操作，去掉数字最右边的 1。
2. 重复此过程，直到数字变为 0，每次去除 1 时计数器加 1。

示例：
对于输入数字 7 (二进制为 0111)，
  - 第一次迭代：
    0111  & (0111 - 1) -> 0110  -> 计数加 1
  - 第二次迭代：
    0110  & (0110 - 1) -> 0100  -> 计数加 1
  - 第三次迭代：
    0100  & (0100 - 1) -> 0000  -> 计数加 1

最终得到计数为 4。

优势：
- 时间复杂度为 O(k)，其中 k 为数字中 1 的数量。
- 显著优于逐位检查所有位的 O(n) 方法，尤其在数字中 1 较少时。

应用：
- 广泛应用于数据压缩、图像处理、网络协议、机器学习等领域，能够提升系统性能。
*/
func CountBitOnePlus(number uint64) int {
	cnt := 0
	for number > 0 {
		// 每次砍掉最右边的1
		number = number & (number - 1)
		cnt++
	}
	return cnt
}
