package cmd

import (
	"os"
	"encoding/json"
	"strings"
	"github.com/Whiteblock/go-prettyjson"
)

func prettyp(s string) string {
	s = strings.Trim(s, "\n\t\r\v ")
	if strings.HasPrefix(s, "[") {
		var temp []interface{}
		err := json.Unmarshal([]byte(s), &temp)
		if err != nil {
			return s
		}
		pps, _ := prettyjson.Marshal(temp)
		return string(pps)
	} else {
		var temp map[string]interface{}

		err := json.Unmarshal([]byte(s), &temp)
		if err != nil {
			return s
		}

		pps, _ := prettyjson.Marshal(temp)

		return string(pps)
	}

}

func prettypi(i interface{}) string {
	_,noPretty := os.LookupEnv("NO_PRETTY")
	if noPretty {
		out,_ := json.Marshal(i)
		return string(out)
	}
	out, _ := prettyjson.Marshal(i)
	return string(out)
}
