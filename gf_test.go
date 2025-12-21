package qrcode

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type GfTestSuite struct {
	suite.Suite
}

func TestGfTestSuite(t *testing.T) {
	suite.Run(t, new(GfTestSuite))
}

func (s *GfTestSuite) TestFindMode() {
	type testcase struct {
		name       string
		polyTarget int
		res        []uint8
	}

	testCases := []testcase{
		{name: "size-20", polyTarget: 20, res: []uint8{0, 17, 60, 79, 50, 61, 163, 26, 187, 202, 180, 221, 225, 83, 239, 156, 164, 212, 212, 188, 190}},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			seedGfTable()
			res, _ := getNextGenPoly(uint8(test.polyTarget))
			s.Equal(res, test.res, "Generator polynomial must be equal")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}
