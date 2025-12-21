package qrcode

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
)

type ErrorCorrection uint8

const (
	L ErrorCorrection = iota
	M
	Q
	H
)

type EncodingMode uint8

const (
	NUMERIC   EncodingMode = 1
	ALPHA     EncodingMode = 2
	BYTE_MODE EncodingMode = 4
	// ECI     EncodingMode = 7
	// KANJI   EncodingMode = 8
)

type QrDetails struct {
	encoding        EncodingMode
	errorCorrection ErrorCorrection
	version         uint8
}

type TableData struct {
}

func FindMode(data []byte) EncodingMode {
	mode := NUMERIC
	for _, item := range data {
		if mode == NUMERIC && slices.Index(numericTable, item) == -1 {
			mode = ALPHA
		}
		if mode == ALPHA && slices.Index(alphaTable, item) == -1 {
			mode = BYTE_MODE
		}
	}
	return mode
}

func FindOptimalVersion(data []byte, correction ErrorCorrection, mode EncodingMode) (int, error) {
	for i := range MAX_VERSION_SUPPORTED {
		v := i + 1
		cap := capacity[v][correction]
		if mode == NUMERIC && int(cap.numeric) < len(data) {
			continue
		}
		if mode == ALPHA && int(cap.alpha) < len(data) {
			continue
		}
		if mode == BYTE_MODE && int(cap.byteMode) < len(data) {
			continue
		}
		config = capacity[v][correction]
		return v, nil
	}
	maxLimit := capacity[MAX_VERSION_SUPPORTED][correction]
	return -1, errors.New(fmt.Sprintf("\nData is too long to support.. \ncheck the limits... \nmax version supported:%d \nnumeric mode:%d \nalpha mode:%d \nbyte mode:%d", MAX_VERSION_SUPPORTED, maxLimit.numeric, maxLimit.alpha, maxLimit.byteMode))
}

func EncodeNumeric(data []byte) []byte {
	newArr := make([]byte, 0, int(len(data)*10/3))

	getBitLen := func(len int) int {
		switch len {
		case 3:
			return 10
		case 2:
			return 7
		case 1:
			return 4
		default:
			return 0
		}
	}

	for i := 0; i < len(data); {
		end := i + 3
		if end > len(data) {
			end = len(data)
		}

		raw := RemoveTrailingZeros(data[i:end], true)
		val, _ := strconv.Atoi(string(raw))

		binArr := NumToBinary(val, getBitLen(len(raw)))
		newArr = append(newArr, binArr...)
		i = end
	}

	return newArr
}

func EncodeAlpha(data []byte) []byte {
	newArr := make([]byte, 0, int(len(data)*11/2))

	for i := 0; i < len(data); {
		end := i + 2
		if end > len(data) {
			end = len(data)
		}

		dt := data[i:end]
		if len(dt) == 2 {
			n1 := slices.Index(alphaTable, dt[0])
			n2 := slices.Index(alphaTable, dt[1])

			binArr := NumToBinary((45*n1)+n2, 11)
			newArr = append(newArr, binArr...)

		} else if len(dt) == 1 {
			n1 := slices.Index(alphaTable, dt[0])

			binArr := NumToBinary(n1, 6)
			newArr = append(newArr, binArr...)

		}

		i = end
	}
	return newArr
}

func EncodeByteMode(data []byte) []byte {
	newArr := make([]byte, 0, int(len(data)*8))

	for _, d := range data {
		binArr := NumToBinary(int(d), 8)
		newArr = append(newArr, binArr...)
	}
	return newArr
}

func EncodeData(data []byte, version int, mode EncodingMode, ec ErrorCorrection) ([]byte, error) {
	loadCharCountFrameSize()
	arr := make([]byte, 0, 100) // adjust later
	arr = append(arr, NumToBinary(int(mode), 4)...)

	switch mode {
	case NUMERIC:
		frameSize := NumToBinary(len(data), int(charCountFrameSize[version].numeric))
		arr = append(arr, frameSize...)

		encoded := EncodeNumeric(data)
		arr = append(arr, encoded...)
		break
	case ALPHA:
		frameSize := NumToBinary(len(data), int(charCountFrameSize[version].alpha))
		arr = append(arr, frameSize...)

		encoded := EncodeAlpha(data)
		arr = append(arr, encoded...)
		break
	case BYTE_MODE:
		frameSize := NumToBinary(len(data), int(charCountFrameSize[version].byteMode))
		arr = append(arr, frameSize...)

		encoded := EncodeByteMode(data)
		arr = append(arr, encoded...)
		break
	default:
		panic("invalid mode")
	}

	wiggleRoom := int(config.totalCodewords*8) - len(arr)
	if wiggleRoom <= 4 {
		arr = append(arr, make([]byte, wiggleRoom)...)
	} else {
		arr = append(arr, make([]byte, 4)...)
		wiggleRoom -= 4

		rem := wiggleRoom % 8
		if rem != 0 {
			arr = append(arr, make([]byte, rem)...)
			wiggleRoom -= rem
		}
		if wiggleRoom > 0 {
			for i := range int(wiggleRoom / 8) {
				if i%2 == 0 {
					arr = append(arr, NumToBinary(236, 8)...)
				} else {
					arr = append(arr, NumToBinary(17, 8)...)
				}
			}
		}
	}
	return arr, nil
}
