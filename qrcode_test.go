package qrcode

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func getReadableBinary(data []byte) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		result[i] = '0' + b
	}
	return result
}

type QrTestSuite struct {
	suite.Suite
}

func TestQrTestSuite(t *testing.T) {
	suite.Run(t, new(QrTestSuite))
}

func (s *QrTestSuite) TestFindMode() {
	type testcase struct {
		name string
		data []byte
		mode EncodingMode
	}

	testCases := []testcase{
		{name: "numeric-test", data: []byte{'1', '2'}, mode: NUMERIC},
		{name: "alpha-test", data: []byte{'1', 'A'}, mode: ALPHA},
		{name: "byteMode-test", data: []byte{'#', 'A'}, mode: BYTE_MODE},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			mode := FindMode(test.data)
			s.Equal(mode, test.mode, "Mode should match")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *QrTestSuite) TestFindOptimalVersion() {
	type testcase struct {
		name       string
		data       []byte
		correction ErrorCorrection
		mode       EncodingMode
		version    int
	}

	testCases := []testcase{
		{name: "numeric-test", data: make([]byte, 27), mode: NUMERIC, correction: Q, version: 1},
		{name: "alpha-test", data: make([]byte, 16), mode: ALPHA, correction: Q, version: 1},
		{name: "byteMode-test", data: make([]byte, 11), mode: BYTE_MODE, correction: Q, version: 1},
		{name: "error", data: make([]byte, 28), mode: NUMERIC, correction: Q, version: -1},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			version, _ := FindOptimalVersion(test.data, test.correction, test.mode)
			s.Equal(version, test.version, "Version should match")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *QrTestSuite) TestEncodeAlpha() {
	type testcase struct {
		name string
		data []byte
		resp string
	}

	testCases := []testcase{
		{
			name: "hello world",
			data: []byte{'H', 'E', 'L', 'L', 'O', ' ', 'W', 'O', 'R', 'L', 'D'},
			resp: "0110000101101111000110100010111001011011100010011010100001101",
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			encoded := EncodeAlpha(test.data)
			s.Equal(string(getReadableBinary(encoded)), test.resp, "Encoded bin should match")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *QrTestSuite) TestEncodeNumeric() {
	type testcase struct {
		name string
		data []byte
		resp string
	}

	testCases := []testcase{
		{
			name: "hello world",
			data: []byte{'8', '6', '7', '5', '3', '0', '9'},
			resp: "110110001110000100101001",
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			encoded := EncodeNumeric(test.data)
			s.Equal(string(getReadableBinary(encoded)), test.resp, "Encoded bin should match")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *QrTestSuite) TestEncodeByteMode() {
	type testcase struct {
		name string
		data []byte
		resp string
	}

	testCases := []testcase{
		{
			name: "hello world",
			data: []byte{'H', 'e', 'l', 'l', 'o'},
			resp: "0100100001100101011011000110110001101111",
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			encoded := EncodeByteMode(test.data)
			s.Equal(string(getReadableBinary(encoded)), test.resp, "Encoded bin should match")
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}

func (s *QrTestSuite) TestEncodeData() {
	type testcase struct {
		name    string
		version int
		mode    EncodingMode
		ec      ErrorCorrection
		data    []byte
		resp    string
	}

	testCases := []testcase{
		{
			name:    "alpha_encode",
			data:    []byte{'H', 'E', 'L', 'L', 'O', ' ', 'W', 'O', 'R', 'L', 'D'},
			version: 1,
			mode:    ALPHA,
			ec:      Q,
			resp:    "00100000010110110000101101111000110100010111001011011100010011010100001101000000111011000001000111101100",
		},
	}

	testingFn := func(test testcase) func(*testing.T) {
		return func(*testing.T) {
			version, _ := FindOptimalVersion(test.data, test.ec, test.mode) // load config variable
			encoded, err := EncodeData(test.data, version, test.mode, config.correction)
			if err != nil {
				s.Equal(len(encoded), 0)
				s.EqualError(err, "invalid mode")
			} else {
				s.Equal(string(getReadableBinary(encoded)), test.resp, "Encoded bin should match")
			}
		}
	}

	t := s.T()

	for _, tt := range testCases {
		t.Run(tt.name, testingFn(tt))
	}
}
