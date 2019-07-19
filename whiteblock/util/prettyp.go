package util

import (
	"encoding/json"
	"fmt"
	"github.com/Whiteblock/go-prettyjson"
	"os"
	"strings"
)

func Prettyp(s string) string {
	_, noPretty := os.LookupEnv("NO_PRETTY")
	if noPretty {
		return s
	}
	s = strings.Trim(s, "\n\t\r\v ")
	if s[0] == '"' {
		return s
	}
	var tmp interface{}
	err := json.Unmarshal([]byte(s), &tmp)
	if err != nil {
		return s
	}
	pps, _ := prettyjson.Marshal(tmp)
	return string(pps)
}

func Prettypi(i interface{}) string {
	_, noPretty := os.LookupEnv("NO_PRETTY")
	if noPretty {
		out, _ := json.Marshal(i)
		return string(out)
	}
	out, _ := prettyjson.Marshal(i)
	return string(out)
}

func Print(i interface{}) {
	switch i.(type) {
	case string:
		_, noPretty := os.LookupEnv("NO_PRETTY")
		if noPretty {
			fmt.Println(i.(string))
		} else {
			fmt.Printf("\033[97m%s\033[0m\n", i.(string))
		}
	default:
		fmt.Println(Prettypi(i))
	}
}

func Printf(format string, a ...interface{}) {
	Print(fmt.Sprintf(format, a...))
}
