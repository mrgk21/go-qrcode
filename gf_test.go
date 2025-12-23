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

func (s *GfTestSuite) TestGetNextGenPoly() {
	type testcase struct {
		name       string
		polyTarget uint8
		res        []uint8
	}

	testCases := []testcase{
		{
			name:       "next-gen-poly",
			polyTarget: 10,
			res:        []uint8{0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
		},
		{
			name:       "next-gen-poly-exists",
			polyTarget: 1,
			res:        []uint8{0, 0},
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			SeedGfTable()
			res, err := GetNextGenPoly(test.polyTarget)
			if err != nil {
				s.EqualError(err, "Error calculating poly")
			}
			s.Equal(res, test.res, "Gen polynomial must be equal")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *GfTestSuite) TestGetMsgPoly() {
	type testcase struct {
		name    string
		msgPoly []uint8
		res     []uint8
	}

	testCases := []testcase{
		{
			name:    "msg-poly",
			msgPoly: []uint8{0, 0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0},
			res:     []uint8{5, 92, 238, 78, 161, 155, 187, 145, 98, 6, 122, 100, 122},
		},
		{
			name:    "msg-poly-error",
			msgPoly: []uint8{0, 1, 0, 0, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 0, 0},
			res:     []uint8{},
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			SeedGfTable()
			res, err := GetMsgPoly(test.msgPoly)
			if err != nil {
				s.EqualError(err, "invalid message length", "message length should be a multiple of 8")
			}
			s.Equal(res, test.res, "Message polynomial must be equal")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *GfTestSuite) TestDividePoly() {
	type testcase struct {
		name string
		msg  []uint8
		gen  []uint8
		op   []uint8
	}

	testCases := []testcase{
		{
			name: "sample1",
			msg:  []uint8{5, 92, 238, 78, 161, 155, 187, 145, 98, 6, 122, 100, 122, 100, 122, 100},
			gen:  []uint8{0, 251, 67, 46, 61, 118, 70, 64, 94, 32, 45},
			op:   []uint8{183, 47, 33, 43, 235, 170, 81, 95, 56, 129},
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			SeedGfTable()
			res := DividePoly(test.msg, test.gen)
			s.Equal(res, test.op, "Remainder polynomial must match")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}
