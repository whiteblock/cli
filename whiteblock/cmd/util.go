package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
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

func GetNodes() ([]Node, error) {
	res, err := jsonRpcCall("nodes", []string{})
	if err != nil {
		return nil, err
	}
	tmp := res.([]interface{})
	nodes := []map[string]interface{}{}
	for _, t := range tmp {
		nodes = append(nodes, t.(map[string]interface{}))
	}

	out := []Node{}
	for _, node := range nodes {
		out = append(out, Node{
			LocalID:   int(node["localId"].(float64)),
			Server:    int(node["server"].(float64)),
			TestNetID: node["testNetId"].(string),
			ID:        int(node["id"].(float64)),
			IP:        node["ip"].(string),
			Label:     node["label"].(string),
		})
	}
	return out, nil
}
