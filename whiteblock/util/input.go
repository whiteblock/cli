package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func YesNoPrompt(msg string) bool {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s ([y]es/no) ", msg)
		scanner.Scan()
		ask := scanner.Text()
		ask = strings.Trim(ask, "\n\t\r\v\f ")

		switch ask {
		case "n":
			fallthrough
		case "no":
			fallthrough
		case "0":
			return false

		case "y":
			fallthrough
		case "yes":
			fallthrough
		case "1":
			return true
		default:
			fmt.Println("Unknown Option")
		}
	}
}
