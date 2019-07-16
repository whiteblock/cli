package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"runtime"
)

const (
	debug     = false
	NoMaxArgs = -1
)

/**
 * Unify error messages through function calls
 */

func CheckArguments(cmd *cobra.Command, args []string, min int, max int) {
	if min == max && len(args) != min {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if min == 1 {
			plural = ""
		}
		PrintStringError(fmt.Sprintf("Invalid number of arguments. Expected exactly %d argument%s. Given %d.", min, plural, len(args)))
		os.Exit(1)
	}
	if len(args) < min {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if min == 1 {
			plural = ""
		}
		PrintStringError(fmt.Sprintf("Missing arguments. Expected atleast %d argument%s. Given %d.", min, plural, len(args)))
		os.Exit(1)
	}
	if max != NoMaxArgs && len(args) > max {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if max == 1 {
			plural = ""
		}
		PrintStringError(fmt.Sprintf("Too many arguments. Expected atmost %d argument%s. Given %d.", max, plural, len(args)))
		os.Exit(1)
	}
}

func InvalidArgument(arg string) {
	PrintStringError(fmt.Sprintf("Invalid argument given: %s.", arg))
}

func InvalidInteger(name string, value string, fatal bool) {
	PrintStringError(fmt.Sprintf("Invalid integer, given \"%s\" for %s.", value, name))
	if fatal {
		os.Exit(1)
	}
}

func CheckIntegerBounds(cmd *cobra.Command, name string, val int, min int, max int) {
	if val < min {
		PrintStringError(fmt.Sprintf("The value given for %s, %d cannot be less than %d.", name, val, min))
		os.Exit(1)
	} else if val > max {
		PrintStringError(fmt.Sprintf("The value given for %s, %d cannot be greater than %d.", name, val, max))
		os.Exit(1)
	}
}

func ClientNotSupported(client string) {
	PrintStringError(fmt.Sprintf("This function is not supported for %s.", client))
	os.Exit(1)
}

func PrintErrorFatal(err error) {
	if debug {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Printf("\033[31mError: %v:%v\033[0m %s\n", file, line, err)
		os.Exit(1)
	}
	PrintError(err)
	os.Exit(1)
}

func PrintError(err error) {
	PrintStringError(err.Error())
}

func PrintStringError(err string) {
	if debug {
		_, file, line, ok := runtime.Caller(1)
		if !ok {
			file = "???"
			line = 0
		}
		fmt.Printf("\033[31mError: %v:%v\033[0m %s\n", file, line, err)
	} else {
		fmt.Printf("\033[31mError:\033[0m %s\n", err)
	}

}
