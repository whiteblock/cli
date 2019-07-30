package util

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestParseIntToStringSlice(t *testing.T) {
	var tests = []struct {
		vals     []string
		expected map[int][]string
	}{
		{
			vals:     []string{"0=blah, blah, blah"},
			expected: map[int][]string{0: []string{"blah, blah, blah"}},
		},
		{
			vals:     []string{"0=blah", "1=test"},
			expected: map[int][]string{0: []string{"blah"}, 1: []string{"test"}},
		},
		{
			vals:     []string{"0=blah, blah, blah", "1=test, test"},
			expected: map[int][]string{0: []string{"blah, blah, blah"}, 1: []string{"test, test"}},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out, err := ParseIntToStringSlice(tt.vals)
			if err != nil {
				t.Error("error running ParseIntToStringSlice", err)
			}

			if !reflect.DeepEqual(out, tt.expected) {
				t.Error("return value of ParseIntToStringSlice does not match expected value")
			}
		})
	}
}

func TestGetAsBool(t *testing.T) {
	var tests = []struct {
		input    string
		expected bool
	}{
		{
			input:    "n",
			expected: false,
		},
		{
			input:    "no",
			expected: false,
		},
		{
			input:    "0",
			expected: false,
		},
		{
			input:    "\r0\f",
			expected: false,
		},
		{
			input:    "y",
			expected: true,
		},
		{
			input:    "yes",
			expected: true,
		},
		{
			input:    "1",
			expected: true,
		},
		{
			input:    "\t1\n",
			expected: true,
		},
		{
			input:    "",
			expected: false,
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			out, _ := GetAsBool(tt.input)

			if out != tt.expected {
				t.Error("return value of GetAsBool does not match expected value")
			}
		})
	}
}

func TestYesNoPrompt(t *testing.T) {
	//var tests = []struct {
	//	msg string
	//	expected bool
	//}{
	//	{
	//		msg: "type y",
	//		expected: true,
	//	},
	//	{
	//		msg: "type yes",
	//		expected: true,
	//	},
	//	{
	//		msg: "type n",
	//		expected: false,
	//	},
	//	{
	//		msg: "type no",
	//		expected: false,
	//	},
	//}
	//
	//for i, tt := range tests {
	//	t.Run(strconv.Itoa(i), func(t *testing.T) {
	//		if YesNoPrompt(tt.msg) != tt.expected {
	//			t.Error("return value of YesNoPrompt does not match expected value")
	//		}
	//	})
	//}

	//TODO fix?
}

func TestArgsToJSON(t *testing.T) {
	args := []string{"blah:somethin", "test:somethin", "num:123", "uInt:0"}


	fmt.Println(ArgsToJSON(args))
}
