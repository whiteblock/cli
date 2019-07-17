package build

import (
	"encoding/json"
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

func SanitizeBuild(conf *Config) {
	conf.Blockchain = strings.ToLower(strings.Trim(conf.Blockchain, "\r\t\v\n "))
	for i := range conf.Images {
		conf.Images[i] = strings.Trim(conf.Images[i], "\r\t\v\n ")
	}
}
