package qrcode

import (
	"errors"
	"math"
)

var gfLogTable = map[int]int{}
var gfAntilogTable = map[int]int{}

func seedGfTable() {
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

func addGf(v1 uint8, v2 uint8) uint8 {
	v := gfLogTable[int(v1)] ^ gfLogTable[int(v2)]
	return uint8(gfAntilogTable[v])
}

func mulGf(v1 int, v2 int) uint8 {
	if (v1 + v2) > 255 {
		return uint8((v1 + v2) % 255)
	}
	return uint8(v1 + v2)
}

var polySeed = []uint8{1, 3, 2}
var maxPolyExp int = 2
var polyMap = map[uint8][]uint8{
	1: []uint8{0, 0},
}

func getNextGenPoly(polyTarget uint8) ([]uint8, error) {
	p, exists := polyMap[polyTarget]
	if exists {
		return p, nil
	}
	p2, exists := polyMap[polyTarget-1]
	if !exists {
		pl2, err := getNextGenPoly(polyTarget - 1)
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
