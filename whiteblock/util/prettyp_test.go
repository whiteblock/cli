package util

import (
	"strconv"
	"testing"
)

func TestPrettyp(t *testing.T) {
	var tests = []struct {
		s string
		expected string
	}{
		{
			s: "some random string\nsome random \rsecond string",
			expected: string([]byte{115, 111, 109, 101, 32, 114, 97, 110, 100, 111, 109, 32, 115, 116, 114, 105, 110, 103, 10, 115, 111, 109, 101, 32, 114, 97, 110, 100, 111, 109, 32, 13, 115, 101, 99, 111, 110, 100, 32, 115, 116, 114, 105, 110, 103}),
		},
		{
			s: "some random string\nsome random \tsecond string",
			expected: string([]byte{115, 111, 109, 101, 32, 114, 97, 110, 100, 111, 109, 32, 115, 116, 114, 105, 110, 103, 10, 115, 111, 109, 101, 32, 114, 97, 110, 100, 111, 109, 32, 9, 115, 101, 99, 111, 110, 100, 32, 115, 116, 114, 105, 110, 103}),

		},
		{
			s: "some random string\nsome random \vsecond string",
			expected: string([]byte{115, 111, 109, 101, 32, 114, 97, 110, 100, 111, 109, 32, 115, 116, 114, 105, 110, 103, 10, 115, 111, 109, 101, 32, 114, 97, 110, 100, 111, 109, 32, 11, 115, 101, 99, 111, 110, 100, 32, 115, 116, 114, 105, 110, 103}),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if Prettyp(tt.s) != tt.expected {
				t.Error("return value of Prettyp does not match expected value")
			}
		})
	}
}

func TestPrettypi(t *testing.T) {
	var tests = []struct {
		i interface{}
		expected string
	}{
		{
			i: []interface{}{123, "123", "blah", false},
			expected: string([]byte{91, 10, 32, 32, 49, 50, 51, 44, 10, 32, 32, 34, 49, 50, 51, 34, 44, 10, 32, 32, 34, 98, 108, 97, 104, 34, 44, 10, 32, 32, 102, 97, 108, 115, 101, 10, 93}),
		},
		{
			i: []interface{}{"1234.01", -34, []byte{1}},
			expected: string([]byte{91, 10, 32, 32, 34, 49, 50, 51, 52, 46, 48, 49, 34, 44, 10, 32, 32, 45, 51, 52, 44, 10, 32, 32, 34, 65, 81, 61, 61, 34, 10, 93}),
		},
		{
			i: []interface{}{byte(0), uint64(1), []string{"i", "f"}},
			expected: string([]byte{91, 10, 32, 32, 48, 44, 10, 32, 32, 49, 44, 10, 32, 32, 91, 10, 32, 32, 32, 32, 34, 105, 34, 44, 10, 32, 32, 32, 32, 34, 102, 34, 10, 32, 32, 93, 10, 93}),
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if Prettypi(tt.i) != tt.expected {
				t.Error("return value of Prettyp does not match expected value")
			}
		})
	}
}


