package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseIntToStringSlice(vals []string) (map[int][]string, error) {
	out := map[int][]string{}
	for _, val := range vals {
		splitVal := strings.SplitN(val, "=", 2)
		if len(splitVal) != 2 {
			return nil, fmt.Errorf("unexpected value %s", val)
		}
		index := CheckAndConvertInt(splitVal[0], "index")
		if _, ok := out[index]; !ok {
			out[index] = []string{splitVal[1]}
		} else {
			out[index] = append(out[index], splitVal[1])
		}
	}
	return out, nil
}

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
		if !scanner.Scan() {
			PrintErrorFatal(scanner.Err())
		}
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
