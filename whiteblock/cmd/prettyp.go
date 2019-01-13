package cmd

import (
	"encoding/json"
	"strings"

	"github.com/hokaccha/go-prettyjson"
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
