package build

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os/user"
	"strings"
	"sync"
)

var _imageTable map[string]map[string]map[string]map[string]string
var _imageTableMux = sync.Mutex{}

func getImageTable() (map[string]map[string]map[string]map[string]string, error) {
	_imageTableMux.Lock()
	defer _imageTableMux.Unlock()
	if _imageTable != nil {
		return _imageTable, nil
	}
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile("/etc/whiteblock.json")
	if err != nil {
		b, err = ioutil.ReadFile(usr.HomeDir + "/cli/etc/whiteblock.json")
		if err != nil {
			b, err = util.HttpRequest("GET", "https://storage.googleapis.com/genesis-public/cli/dev/etc/whiteblock.json", "")
			if err != nil {
				return nil, err
			}
		}
	}
	err = json.Unmarshal(b, &_imageTable)
	if err != nil {
		return nil, err
	}
	return _imageTable, nil
}

func determineImage(blockchain string, requested string) string {

	cont, err := getImageTable()
	if err != nil {
		if len(requested) > 0 {
			return requested
		}
		util.PrintErrorFatal(err)
	}

	defaultImage := "gcr.io/whiteblock/" + blockchain + ":master"

	log.WithFields(log.Fields{"imageTable": cont}).Trace("parsed the image table")

	if _, ok := cont["blockchains"][blockchain]; !ok {
		log.Debug("chose default image due to missing entry")
		if len(requested) > 0 {
			return requested
		}
		return defaultImage
	}
	if _, ok := cont["blockchains"][blockchain]["images"]; !ok {
		log.Debug("chose default image due to missing entry")
		if len(requested) > 0 {
			return requested
		}
		return defaultImage
	}

	if len(requested) == 0 {
		if stableImage, ok := cont["blockchains"][blockchain]["images"]["stable"]; ok {
			return stableImage
		}
		log.Debugf("missing default stable image for %s", blockchain)
		return defaultImage
	}

	if image, ok := cont["blockchains"][blockchain]["images"][requested]; ok {
		return image
	}
	return requested
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
	return util.CheckAndConvertInt(in[:index], "node"), in[index:len(in)]
}
