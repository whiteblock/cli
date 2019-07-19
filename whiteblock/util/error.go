package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

const (
	NoMaxArgs          = -1
	ErrNoPreviousBuild = "No previous build found"
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
		PrintErrorFatal(fmt.Sprintf("Invalid number of arguments. Expected exactly %d argument%s. Given %d.", min, plural, len(args)))
	}
	if len(args) < min {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if min == 1 {
			plural = ""
		}
		PrintErrorFatal(fmt.Sprintf("Missing arguments. Expected atleast %d argument%s. Given %d.", min, plural, len(args)))
	}
	if max != NoMaxArgs && len(args) > max {
		fmt.Println(cmd.UsageString())
		plural := "s"
		if max == 1 {
			plural = ""
		}
		PrintErrorFatal(fmt.Sprintf("Too many arguments. Expected atmost %d argument%s. Given %d.", max, plural, len(args)))
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
	PrintErrorFatal(fmt.Sprintf("This function is not supported for %s.", client))
}

func PrintErrorFatal(err interface{}) {
	PrintStringError(fmt.Sprint(err))
	os.Exit(1)
}

func PrintStringError(err string) {
	fmt.Printf("\033[31mError:\033[0m %s\n", err)
}
