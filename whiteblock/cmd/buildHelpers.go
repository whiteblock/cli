package cmd

import (
	util "../util"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"os/user"
	"strconv"
	"strings"
)

func getPreviousBuildId() (string, error) {
	buildId, err := util.ReadStore(".previous_build_id")
	if err != nil || len(buildId) == 0 {
		return "", fmt.Errorf("No previous build. Use build command to deploy a blockchain.")
	}
	return string(buildId), nil
}

func getPreviousBuild() (Config, error) {
	buildId, err := getPreviousBuildId()
	if err != nil {
		return Config{}, err
	}

	prevBuild, err := jsonRpcCall("get_build", []string{buildId})
	if err != nil {
		return Config{}, err
	}

	tmp, err := json.Marshal(prevBuild)
	if err != nil {
		return Config{}, err
	}

	var out Config
	err = json.Unmarshal(tmp, &out)
	return out, err
}

func fetchPreviousBuild() (Config, error) {
	buildId, err := getPreviousBuildId()
	if err != nil {
		return Config{}, err
	}

	prevBuild, err := jsonRpcCall("get_last_build", []string{buildId})
	if err != nil {
		return Config{}, err
	}

	tmp, err := json.Marshal(prevBuild)
	if err != nil {
		return Config{}, err
	}

	var out Config
	err = json.Unmarshal(tmp, &out)
	return out, err
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
	rawOptions, err := jsonRpcCall("get_params", []string{blockchain})
	if err != nil {
		return nil, err
	}
	optionsStep1, ok := rawOptions.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Unexpected format for params")
	}
	out := make([][]string, len(optionsStep1))
	for i, optionsStep1Segment := range optionsStep1 {
		optionsStep2, ok := optionsStep1Segment.([]interface{}) //[][]interface{}[i]
		if !ok {
			return nil, fmt.Errorf("Unexpected format for params")
		}
		out[i] = make([]string, len(optionsStep2))
		for j, optionsStep2Segment := range optionsStep2 {
			out[i][j], ok = optionsStep2Segment.(string)
			if !ok {
				return nil, fmt.Errorf("Unexpected format for params")
			}
		}
	}
	return out, nil
}

func getServer() []int {
	idList := make([]int, 0)
	res, err := jsonRpcCall("get_servers", []string{})
	if err != nil {
		util.PrintErrorFatal(err)
	}
	servers := res.(map[string]interface{})
	serverID := 0
	for _, v := range servers {
		serverID = int(v.(map[string]interface{})["id"].(float64))
		//move this and take out break statement if instance has multiple servers
		idList = append(idList, serverID)
		break
	}

	return idList
}

func tern(exp bool, res1 string, res2 string) string {
	if exp {
		return res1
	}
	return res2
}

func removeSmartContracts() {
	cwd := os.Getenv("HOME")
	err := os.RemoveAll(cwd + "/smart-contracts/whiteblock/contracts.json")
	if err != nil {
		util.PrintErrorFatal(err)
	}
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
		util.PrintStringError("Cannot have a numerical environment variable")
		os.Exit(1)
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

func handleDockerAuthFlags(cmd *cobra.Command, args []string, conf *Config) {
	if cmd.Flags().Changed("docker-password") != cmd.Flags().Changed("docker-username") {
		if cmd.Flags().Changed("docker-password") {
			util.PrintStringError("You must also provide --docker-password with --docker-username")
		} else {
			util.PrintStringError("You must also provide --docker-username with --docker-password")
		}
		os.Exit(1)
	}
	if !cmd.Flags().Changed("docker-password") {
		return //The auth flags have not been set
	}

	_, ok := conf.Extras["prebuild"]
	if !ok {
		conf.Extras["prebuild"] = map[string]interface{}{}
	}
	username, err := cmd.Flags().GetString("docker-username")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	password, err := cmd.Flags().GetString("docker-password")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	conf.Extras["prebuild"].(map[string]interface{})["auth"] = map[string]string{
		"username": username,
		"password": password,
	}

}

func handlePullFlag(cmd *cobra.Command, args []string, conf *Config) {
	_, ok := conf.Extras["prebuild"]
	if !ok {
		conf.Extras["prebuild"] = map[string]interface{}{}
	}
	fbg, err := cmd.Flags().GetBool("force-docker-pull")
	if err == nil && fbg {
		conf.Extras["prebuild"].(map[string]interface{})["pull"] = true
	}
}

func handleForceUnlockFlag(cmd *cobra.Command, args []string, conf *Config) {

	fbg, err := cmd.Flags().GetBool("force-unlock")
	if err == nil && fbg {
		conf.Extras["forceUnlock"] = true
	}
}

func getImage(blockchain string, imageType string, defaultImage string) string {
	usr, err := user.Current()
	if err != nil {
		util.PrintErrorFatal(err)
	}
	b, err := ioutil.ReadFile("/etc/whiteblock.json")
	if err != nil {
		b, err = ioutil.ReadFile(usr.HomeDir + "/cli/etc/whiteblock.json")
		if err != nil {
			b, err = util.HttpRequest("GET", "https://whiteblock.io/releases/cli/v1.5.7/whiteblock.json", "")
			if err != nil {
				util.PrintErrorFatal(err)
			}
		}
	}
	var cont map[string]map[string]map[string]map[string]string
	err = json.Unmarshal(b, &cont)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	//fmt.Printf("%#v\n",cont["blockchains"])
	if len(defaultImage) > 0 {
		return defaultImage
	} else if len(cont["blockchains"][blockchain]["images"][imageType]) != 0 {
		return cont["blockchains"][blockchain]["images"][imageType]
	} else {
		return "gcr.io/whiteblock/" + blockchain + ":master"
	}
}

func handleImageFlag(cmd *cobra.Command, args []string, conf *Config) {

	imageFlag, err := cmd.Flags().GetStringSlice("image")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	conf.Images = make([]string, conf.Nodes)
	images, potentialImage, err := util.UnrollStringSliceToMapIntString(imageFlag, "=")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	//fmt.Printf("IMAGES=%#v\n",images)
	if len(potentialImage) > 1 {
		util.PrintErrorFatal(fmt.Errorf("Too many default images"))
	}
	imgDefault := ""
	if len(potentialImage) == 1 {
		fmt.Println("given default image")
		imgDefault = potentialImage[0]
	}
	baseImage := getImage(conf.Blockchain, "stable", imgDefault)

	for i := 0; i < conf.Nodes; i++ {

		conf.Images[i] = baseImage
		image, exists := images[i]
		if exists {
			fmt.Println("exists")
			conf.Images[i] = image
		}
		//fmt.Println(conf.Images[i])
	}
}

func handleFilesFlag(cmd *cobra.Command, args []string, conf *Config) {
	filesFlag, err := cmd.Flags().GetStringSlice("template")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	if filesFlag == nil {
		return
	}

	conf.Files = make([]map[string]string, conf.Nodes)
	defaults := map[string]string{}
	for _, tfileIn := range filesFlag {
		tuple := strings.SplitN(tfileIn, ";", 3) //support both delim in future
		if len(tuple) < 3 {
			tmp := strings.Replace(tfileIn, ";", "=", 1)
			tuple = strings.SplitN(tmp, "=", 2)
			if len(tuple) != 2 {
				util.PrintErrorFatal(fmt.Errorf("Invalid argument"))
			}
		}
		if len(tuple) == 2 {
			data, err := ioutil.ReadFile(tuple[1])
			if err != nil {
				util.PrintErrorFatal(err)
			}
			defaults[tuple[0]] = base64.StdEncoding.EncodeToString(data)
			continue
		}
		data, err := ioutil.ReadFile(tuple[2])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		index, err := strconv.Atoi(tuple[0])
		if err != nil {
			util.PrintErrorFatal(err)
		}
		if index < 0 || index >= conf.Nodes {
			util.PrintErrorFatal(fmt.Errorf("Index is out of range for -t flag"))
		}
		conf.Files[index] = map[string]string{}
		conf.Files[index][tuple[1]] = base64.StdEncoding.EncodeToString(data)
	}

	if conf.Extras == nil {
		conf.Extras = map[string]interface{}{}
	}
	if _, ok := conf.Extras["defaults"]; !ok {
		conf.Extras["defaults"] = map[string]interface{}{}
	}
	conf.Extras["defaults"].(map[string]interface{})["files"] = defaults
}

func handleSSHOptions(cmd *cobra.Command, args []string, conf *Config) {
	if !cmd.Flags().Changed("user-ssh-key") { //Don't bother if not specified
		return
	}

	sshPubKeys, err := cmd.Flags().GetStringSlice("user-ssh-key")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	if conf.Extras == nil {
		conf.Extras = map[string]interface{}{}
	}
	if _, ok := conf.Extras["postbuild"]; !ok {
		conf.Extras["postbuild"] = map[string]interface{}{}
	}
	if _, ok := conf.Extras["postbuild"].(map[string]interface{})["ssh"]; !ok {
		conf.Extras["postbuild"].(map[string]interface{})["ssh"] = map[string]interface{}{}
	}
	pubKeys := []string{}
	for _, pubKeyFile := range sshPubKeys {
		data, err := ioutil.ReadFile(pubKeyFile)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		pubKeys = append(pubKeys, string(data))
	}

	conf.Extras["postbuild"].(map[string]interface{})["ssh"].(map[string]interface{})["pubKeys"] = pubKeys
}
