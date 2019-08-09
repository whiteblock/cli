package build

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os/user"
	"strings"
)

func getImage(blockchain string, imageType string, defaultImage string) string {
	usr, err := user.Current()
	if err != nil {
		util.PrintErrorFatal(err)
	}
	b, err := ioutil.ReadFile("/etc/whiteblock.json")
	if err != nil {
		b, err = ioutil.ReadFile(usr.HomeDir + "/cli/etc/whiteblock.json")
		if err != nil {
			b, err = util.HttpRequest("GET", "https://storage.googleapis.com/genesis-public/cli/dev/etc/whiteblock.json", "")
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

func SanitizeBuild(conf *Config) {
	conf.Blockchain = strings.ToLower(strings.Trim(conf.Blockchain, "\r\t\v\n "))
	for i := range conf.Images {
		conf.Images[i] = strings.Trim(conf.Images[i], "\r\t\v\n ")
	}
}

func getServer() []int {
	idList := make([]int, 0)
	res, err := util.JsonRpcCall("get_servers", []string{})
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

func GetPreviousBuildIDErr() (string, error) {
	var buildID string
	err := util.GetP("previous_build_id", &buildID)
	if err != nil || len(buildID) == 0 {
		return "", fmt.Errorf("No previous build. Use build command to deploy a blockchain, " +
			"or run `whiteblock sync` if you already have a blockchain deployed.")
	}
	return buildID, nil
}

//"github.com/sirupsen/logrus"
func GetPreviousBuildID() string {
	res, err := GetPreviousBuildIDErr()
	if err != nil {
		util.PrintErrorFatal(err)
	}
	return res
}

func GetPreviousBuild() (Config, error) {
	buildId, err := GetPreviousBuildIDErr()
	if err != nil {
		return Config{}, err
	}

	prevBuild, err := util.JsonRpcCall("get_build", []string{buildId})
	if err != nil {
		return Config{}, err
	}
	log.WithFields(log.Fields{"fetched": prevBuild}).Debug("fetched the previous build")

	tmp, err := json.Marshal(prevBuild)
	if err != nil {
		return Config{}, err
	}

	var out Config
	return out, json.Unmarshal(tmp, &out)
}

func FetchPreviousBuild() (Config, error) {
	buildId, err := GetPreviousBuildIDErr()
	if err != nil {
		return Config{}, err
	}

	prevBuild, err := util.JsonRpcCall("get_last_build", []string{buildId})
	if err != nil {
		return Config{}, err
	}

	tmp, err := json.Marshal(prevBuild)
	if err != nil {
		return Config{}, err
	}

	var out Config
	return out, json.Unmarshal(tmp, &out)
}
