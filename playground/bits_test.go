package playground

import "testing"

func TestIsBitOne(t *testing.T) {
	tests := []struct {
		number      uint64
		bitPosition uint
		expected    bool
		shouldPanic bool
	}{
		{number: 5, bitPosition: 1, expected: true, shouldPanic: false},  // 5 = 0101, 第1位是1
		{number: 5, bitPosition: 2, expected: false, shouldPanic: false}, // 5 = 0101, 第2位是0
		{number: 5, bitPosition: 3, expected: true, shouldPanic: false},  // 5 = 0101, 第3位是1
		{number: 5, bitPosition: 65, shouldPanic: true},                  // 超出范围，应该触发 panic
	}

	for _, tt := range tests {
		if tt.shouldPanic {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for bitPosition %d, but did not panic", tt.bitPosition)
				}
			}()
			IsBitOne(tt.number, tt.bitPosition)
		} else {
			actual := IsBitOne(tt.number, tt.bitPosition)
			if actual != tt.expected {
				t.Errorf("IsBitOne(%d, %d) = %v; want %v", tt.number, tt.bitPosition, actual, tt.expected)
			}
		}
	}
}

func TestSetBitOne(t *testing.T) {
	tests := []struct {
		number      uint64
		bitPosition uint
		expected    uint64
		shouldPanic bool
	}{
		{number: 5, bitPosition: 2, expected: 7, shouldPanic: false},  // 5 = 0101, 设置第2位为1 -> 0111 = 7
		{number: 0, bitPosition: 1, expected: 1, shouldPanic: false},  // 0 = 0000, 设置第1位为1 -> 0001 = 1
		{number: 7, bitPosition: 4, expected: 15, shouldPanic: false}, // 7 = 0111, 设置第4位为1 -> 1111 = 15
		{number: 5, bitPosition: 65, shouldPanic: true},               // 超出范围，应该触发 panic
	}

	for _, tt := range tests {
		if tt.shouldPanic {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for bitPosition %d, but did not panic", tt.bitPosition)
				}
			}()
			UpdateBitToOne(tt.number, tt.bitPosition)
		} else {
			number := tt.number
			UpdateBitToOne(number, tt.bitPosition)
			if UpdateBitToOne(number, tt.bitPosition) != tt.expected {
				t.Errorf("UpdateBitToOne(%d, %d) = %d; want %d", tt.number, tt.bitPosition, number, tt.expected)
			}
		}
	}
}

func TestCountBitOne(t *testing.T) {
	tests := []struct {
		number   uint64
		expected int
	}{
		{number: 0, expected: 0},                     // 0 = 0000, 没有1
		{number: 1, expected: 1},                     // 1 = 0001, 有1个1
		{number: 7, expected: 3},                     // 7 = 0111, 有3个1
		{number: 15, expected: 4},                    // 15 = 1111, 有4个1
		{number: 1023, expected: 10},                 // 1023 = 1111111111, 有10个1
		{number: 18446744073709551615, expected: 64}, // 全1（64位）
	}

	for _, tt := range tests {
		actual := CountBitOne(tt.number)
		if actual != tt.expected {
			t.Errorf("CountBitOne(%d) = %d; want %d", tt.number, actual, tt.expected)
		}
	}
}

func TestCountBitOnePlus(t *testing.T) {
	tests := []struct {
		number   uint64
		expected int
	}{
		{number: 0, expected: 0},                     // 0 = 0000, 没有1
		{number: 1, expected: 1},                     // 1 = 0001, 有1个1
		{number: 7, expected: 3},                     // 7 = 0111, 有3个1
		{number: 15, expected: 4},                    // 15 = 1111, 有4个1
		{number: 1023, expected: 10},                 // 1023 = 1111111111, 有10个1
		{number: 18446744073709551615, expected: 64}, // 全1（64位）
	}

	for _, tt := range tests {
		actual := CountBitOnePlus(tt.number)
		if actual != tt.expected {
			t.Errorf("CountBitOnePlus(%d) = %d; want %d", tt.number, actual, tt.expected)
		}
	}
}
