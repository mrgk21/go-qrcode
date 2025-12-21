package qrcode

// modes
var numericTable = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
var alphaTable = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', ' ', '$', '%', '*', '+', '-', '.', '/', ':'}

// char capacity
type charCapacity struct {
	correction     ErrorCorrection
	numeric        uint16
	alpha          uint16
	byteMode       uint16
	totalCodewords uint16
}

const MAX_VERSION_SUPPORTED = 1
const MAX_VERSION_QR = 40

// key = version no.
var capacity = map[int]map[ErrorCorrection]*charCapacity{
	1: {
		L: &charCapacity{correction: L, numeric: 41, alpha: 25, byteMode: 17, totalCodewords: 19},
		M: &charCapacity{correction: M, numeric: 34, alpha: 20, byteMode: 14, totalCodewords: 16},
		Q: &charCapacity{correction: Q, numeric: 27, alpha: 16, byteMode: 11, totalCodewords: 13},
		H: &charCapacity{correction: H, numeric: 17, alpha: 10, byteMode: 7, totalCodewords: 9},
	},
}
var config *charCapacity

type charFrameSize struct {
	numeric  uint8
	alpha    uint8
	byteMode uint8
}

var charCountFrameSize = map[int]*charFrameSize{}

func loadCharCountFrameSize() {
	for i := range MAX_VERSION_SUPPORTED {
		v := i + 1
		if v >= 27 {
			charCountFrameSize[v] = &charFrameSize{numeric: 14, alpha: 13, byteMode: 16}
			continue
		}
		if v >= 10 {
			charCountFrameSize[v] = &charFrameSize{numeric: 12, alpha: 11, byteMode: 16}
			continue
		}
		if v >= 1 {
			charCountFrameSize[v] = &charFrameSize{numeric: 10, alpha: 9, byteMode: 8}
			continue
		}
	}
}
