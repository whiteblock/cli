package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
	"strconv"
	"strings"
)

//"github.com/sirupsen/logrus"
func getPreviousBuildId() (string, error) {
	var buildID string
	err := util.GetP("previous_build_id", &buildID)
	if err != nil || len(buildID) == 0 {
		return "", fmt.Errorf("No previous build. Use build command to deploy a blockchain, " +
			"or run `whiteblock sync` if you already have a blockchain deployed.")
	}
	return buildID, nil
}

func getPreviousBuild() (build.Config, error) {
	buildId, err := getPreviousBuildId()
	if err != nil {
		return build.Config{}, err
	}

	prevBuild, err := util.JsonRpcCall("get_build", []string{buildId})
	if err != nil {
		return build.Config{}, err
	}
	log.WithFields(log.Fields{"fetched": prevBuild}).Debug("fetched the previous build")

	tmp, err := json.Marshal(prevBuild)
	if err != nil {
		return build.Config{}, err
	}

	var out build.Config
	return out, json.Unmarshal(tmp, &out)
}

func fetchPreviousBuild() (build.Config, error) {
	buildId, err := getPreviousBuildId()
	if err != nil {
		return build.Config{}, err
	}

	prevBuild, err := util.JsonRpcCall("get_last_build", []string{buildId})
	if err != nil {
		return build.Config{}, err
	}

	tmp, err := json.Marshal(prevBuild)
	if err != nil {
		return build.Config{}, err
	}

	var out build.Config
	return out, json.Unmarshal(tmp, &out)
}

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

func processOptions(givenOptions map[string]string, format [][]string) (map[string]interface{}, error) {
	out := map[string]interface{}{}

	for _, kv := range format {
		name := kv[0]
		key_type := kv[1]

		val, ok := givenOptions[name]
		if !ok {
			continue
		}
		switch key_type {
		case "string":
			//needs to have filtering
			out[name] = val
		case "[]string":
			preprocessed := strings.Replace(val, " ", ",", -1)
			out[name] = strings.Split(preprocessed, ",")
		case "int":
			val, err := strconv.ParseInt(val, 0, 64)
			if err != nil {
				return nil, err
			}
			out[name] = val

		case "bool":
			switch val {
			case "true":
				fallthrough
			case "yes":
				out[name] = true
			case "false":
				fallthrough
			case "no":
				out[name] = false
			}
		}
	}
	return out, nil
}

//-1 means for all
func processEnvKey(in string) (int, string) {
	node := -1
	index := 0
	for i, char := range in {
		if char < '0' || char > '9' {
			index = i
			break
		}
	}
	if index == 0 {
		return node, in
	}

	if index == len(in) {
		util.PrintErrorFatal("Cannot have a numerical environment variable")
	}

	var err error
	node, err = strconv.Atoi(in[:index])
	if err != nil {
		util.PrintErrorFatal(err)
	}
	return node, in[index:len(in)]
}

func processEnv(envVars map[string]string, nodes int) ([]map[string]string, error) {
	out := make([]map[string]string, nodes)
	for i, _ := range out {
		out[i] = make(map[string]string)
	}
	for k, v := range envVars {
		node, key := processEnvKey(k)
		if node == -1 {
			for i, _ := range out {
				out[i][key] = v
			}
			continue
		}
		out[node][key] = v
	}
	return out, nil
}
