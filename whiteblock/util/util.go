package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
)

func PartialCommand(cmd *cobra.Command, args []string) {
	fmt.Println("\nNo command given. Please choose a command from the list below.")
	cmd.Help()
	return
}

func CheckAndConvertInt(num string, name string) int {
	out, err := strconv.ParseInt(num, 0, 32)
	if err != nil {
		InvalidInteger(name, num, true)
	}
	return int(out)
}

func CheckAndConvertInt64(num string, name string) int64 {
	out, err := strconv.ParseInt(num, 0, 64)
	if err != nil {
		InvalidInteger(name, num, true)
	}
	return out
}

/*
   Write writes data to a file, creating it if it doesn't exist,
   deleting and recreating it if it does.
*/
func Write(path string, data []byte) error {
	return ioutil.WriteFile(path, data, 0664)
}
