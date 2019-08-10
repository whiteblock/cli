package cmd

import (
	"fmt"
	"github.com/whiteblock/cli/whiteblock/util"
)

func hasParam(params [][]string, param string) bool {
	for _, p := range params {
		if p[0] == param {
			return true
		}
	}
	return false
}

func fetchParams(blockchain string) ([][]string, error) {
	//Handle the ugly conversions, in a safe manner
	badFmtErr := fmt.Errorf("unexpected format for params")
	rawOptions, err := util.JsonRpcCall("get_params", []string{blockchain})
	if err != nil {
		return nil, err
	}
	optionsStep1, ok := rawOptions.([]interface{})
	if !ok {
		return nil, badFmtErr
	}
	out := make([][]string, len(optionsStep1))
	for i, optionsStep1Segment := range optionsStep1 {
		optionsStep2, ok := optionsStep1Segment.([]interface{}) //[][]interface{}[i]
		if !ok {
			return nil, badFmtErr
		}
		out[i] = make([]string, len(optionsStep2))
		for j, optionsStep2Segment := range optionsStep2 {
			out[i][j], ok = optionsStep2Segment.(string)
			if !ok {
				return nil, badFmtErr
			}
		}
	}
	return out, nil
}

func tern(exp bool, res1 string, res2 string) string {
	if exp {
		return res1
	}
	return res2
}
