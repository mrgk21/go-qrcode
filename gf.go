package qrcode

import (
	"errors"
	"math"
)

var gfLogTable = map[int]int{}
var gfAntilogTable = map[int]int{}

func SeedGfTable() {
	gfLogTable[0] = 1
	for i := range 255 {
		power := i + 1
		if power < 8 {
			count := int(math.Pow(2, float64(power)))
			gfLogTable[power] = count
			gfAntilogTable[count] = power
			continue
		}
		prevVal := gfLogTable[power-1]
		if prevVal*2 <= 255 {
			gfLogTable[power] = prevVal * 2
			gfAntilogTable[prevVal*2] = power
		} else {
			count := prevVal*2 ^ 285
			gfLogTable[power] = count
			gfAntilogTable[count] = power
		}
	}
}

// ip -> alpha, op -> alpha
func addGf(v1 uint8, v2 uint8) uint8 {
	v := gfLogTable[int(v1)] ^ gfLogTable[int(v2)]
	return uint8(gfAntilogTable[v])
}

// ip -> alpha, op -> alpha
func mulGf(v1 int, v2 int) uint8 {
	if (v1 + v2) > 255 {
		return uint8((v1 + v2) % 255)
	}
	return uint8(v1 + v2)
}

var polyMap = map[uint8][]uint8{
	1: []uint8{0, 0},
}

func GetNextGenPoly(polyTarget uint8) ([]uint8, error) {
	p, exists := polyMap[polyTarget]
	if exists {
		return p, nil
	}
	p2, exists := polyMap[polyTarget-1]
	if !exists {
		pl2, err := GetNextGenPoly(polyTarget - 1)
		if err != nil {
			return []uint8{}, errors.New("Error calculating poly")
		}
		p2 = pl2
	}
	polyArr := make([]uint8, 0, polyTarget)
	m := []uint8{0, polyTarget - 1}

	for i := range polyTarget + 1 {
		if i == 0 {
			polyArr = append(polyArr, mulGf(int(p2[i]), int(m[0])))
		} else if i == polyTarget {
			polyArr = append(polyArr, mulGf(int(p2[polyTarget-1]), int(m[1])))
		} else {
			polyArr = append(polyArr, addGf(mulGf(int(p2[i]), 0), mulGf(int(p2[i-1]), int(m[1]))))
		}
	}
	polyMap[polyTarget] = polyArr
	return polyArr, nil
}

func GetMsgPoly(msg []uint8) ([]uint8, error) {
	if (len(msg) % 8) != 0 {
		return []uint8{}, errors.New("invalid message length")
	}

	val := make([]uint8, 0, int(len(msg)/8))
	for i := 0; i < len(msg); i += 8 {
		val = append(val, uint8(gfAntilogTable[BinToNum(msg[i:i+8])]))
	}
	return val, nil
}

func PrintQ(msg []uint8) []uint8 {
	msg2 := make([]uint8, len(msg))
	for i, v := range msg {
		msg2[i] = uint8(gfLogTable[int(v)])
	}
	return msg2
}

// msg/gen
func DividePoly(msg []uint8, gen []uint8) []uint8 {
	msg2 := make([]uint8, len(msg)+len(gen)-1)
	copy(msg2, msg)

	originalLen := len(msg)
	for i, msgItem := range msg2 {
		if i == originalLen {
			break
		}
		for i1, genTerm := range gen {
			adjustedGen := mulGf(int(genTerm), int(msgItem))
			if msg2[i1] == 0 {
				msg2[i1] = adjustedGen
			} else {
				msg2[i1] = addGf(msg2[i1], adjustedGen)
			}
		}
		msg2 = msg2[1:]
	}

	return msg2
}
