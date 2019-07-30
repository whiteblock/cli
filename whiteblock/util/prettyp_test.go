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


