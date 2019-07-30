package util

import (
	"strconv"
	"testing"
)

func TestPartialCommand(t *testing.T) {
	//var tests = []struct {
	//	cmd *cobra.Command
	//	args []string
	//}{
	//	{
	//		cmd: &cobra.Command{},
	//		args: []string{},
	//	},
	//	{
	//		cmd: &cobra.Command{},
	//		args: []string{"one"},
	//	},
	//	{
	//		cmd: &cobra.Command{},
	//		args: []string{"one", "two"},
	//	},
	//}

	// How to compare prints?

}

func TestCheckAndConvertInt_Successful(t *testing.T) {
	var tests = []struct {
		num string
		name string
		expected int
	}{
		{num: "5", name: "test", expected: 5},
		{num: "158348", name: "test", expected: 158348},
		{num: "0", name: "test", expected: 0},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			if CheckAndConvertInt(tt.num, tt.name) != tt.expected {
				t.Error("return value of CheckAndConvertInt does not match expected value")
			}
		})
	}
}