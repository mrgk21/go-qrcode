package qrcode

import (
	"slices"
)

func RemoveTrailingZeros(arr []byte, isCharFormat bool) []byte {
	var compByte uint8 = 0
	if isCharFormat {
		compByte = '0'
	}
	for i, b := range arr {
		if b != compByte {
			return arr[i:]
		}
	}
	return arr
}

// merge with numToBinary
func CharToBinary(data byte, bitSize int) []byte {
	newArr := make([]byte, 0, bitSize)
	for j := range bitSize {
		zeroOrOne := data >> (7 - j) & 1
		newArr = append(newArr, zeroOrOne)
	}
	return newArr
}

func NumToBinary(num int, bitSize int) []byte {
	newArr := make([]byte, 0, bitSize)
	for true {
		rem := num & 1
		num = (num >> 1)
		newArr = append(newArr, byte(rem))
		if num == 0 {
			break
		}
	}

	newArr = append(newArr, make([]byte, cap(newArr)-len(newArr))...)
	slices.Reverse(newArr)
	return newArr
}

func BinToNum(num []uint8) int {
	var val int = 0
	for _, v := range num {
		val = val*2 + int(v)
	}
	return val
}

func And()        {}
func Or()         {}
func Xor()        {}
func Append()     {}
func PadTo8Bits() {}
