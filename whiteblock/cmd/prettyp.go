package cmd

import (
	"encoding/json"

	"github.com/hokaccha/go-prettyjson"
)

func prettyp(s string) string {
	var temp map[string]interface{}

	err := json.Unmarshal([]byte(s), &temp)
	if err != nil {
		panic(err)
	}

	pps, _ := prettyjson.Marshal(temp)

	return string(pps)
}
