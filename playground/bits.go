package playground

// IsBitOne 判断第i位上是否为1，i从1开始
func IsBitOne(number uint64, bitPosition uint) bool {
	if bitPosition < 1 || bitPosition > 64 {
		panic("index out of range")
	}
	return number&(1<<(bitPosition-1)) != 0
}
