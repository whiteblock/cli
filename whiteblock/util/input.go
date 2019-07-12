package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetAsBool(input string) (bool, error) {
	switch strings.Trim(input, "\n\t\r\v\f ") {
	case "n":
		fallthrough
	case "no":
		fallthrough
	case "0":
		return false, nil

	case "y":
		fallthrough
	case "yes":
		fallthrough
	case "1":
		return true, nil
	}
	return false, fmt.Errorf("Unknown option for boolean")
}

func YesNoPrompt(msg string) bool {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("%s ([y]es/[n]o) ", msg)
		scanner.Scan()
		ask := scanner.Text()
		res, err := GetAsBool(ask)
		if err != nil {
			fmt.Println(err)
			continue
		}
		return res
	}
	panic("should never reach")
}
