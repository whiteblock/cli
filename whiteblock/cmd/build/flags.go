package build

import (
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func HandleForceUnlockFlag(cmd *cobra.Command, args []string, conf *Config) {

	fbg, err := cmd.Flags().GetBool("force-unlock")
	if err == nil && fbg {
		conf.Extras["forceUnlock"] = true
	}
}

func HandlePullFlag(cmd *cobra.Command, args []string, conf *Config) {
	_, ok := conf.Extras["prebuild"]
	if !ok {
		conf.Extras["prebuild"] = map[string]interface{}{}
	}
	fbg, err := cmd.Flags().GetBool("force-docker-pull")
	if err == nil && fbg {
		conf.Extras["prebuild"].(map[string]interface{})["pull"] = true
	}
}

func HandleDockerAuthFlags(cmd *cobra.Command, args []string, conf *Config) {
	if cmd.Flags().Changed("docker-password") != cmd.Flags().Changed("docker-username") {
		if cmd.Flags().Changed("docker-password") {
			util.PrintStringError("you must also provide --docker-password with --docker-username")
		} else {
			util.PrintStringError("you must also provide --docker-username with --docker-password")
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

func HandleImageFlag(cmd *cobra.Command, args []string, conf *Config) {

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
		imgDefault = potentialImage[0]
		log.WithFields(log.Fields{"image": imgDefault}).Debug("given default image")
	}
	baseImage := getImage(conf.Blockchain, "stable", imgDefault)

	for i := 0; i < conf.Nodes; i++ {

		conf.Images[i] = baseImage
		image, exists := images[i]
		if exists {
			log.WithFields(log.Fields{"image": image}).Trace("image exists")
			conf.Images[i] = image
		}
	}
}

func HandleFilesFlag(cmd *cobra.Command, args []string, conf *Config) {
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
		for i := range tuple {
			tuple[i] = strings.Trim(tuple[i], " \n\r\t")
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

func HandleSSHOptions(cmd *cobra.Command, args []string, conf *Config) {
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

func HandleDockerfile(cmd *cobra.Command, args []string, conf *Config) {
	filePath, err := cmd.Flags().GetString("dockerfile")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	if len(filePath) == 0 {
		return
	}
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	if conf.Extras == nil {
		conf.Extras = map[string]interface{}{}
	}

	if _, ok := conf.Extras["prebuild"]; !ok {
		conf.Extras["prebuild"] = map[string]interface{}{}
	}
	conf.Extras["prebuild"].(map[string]interface{})["build"] = true
	conf.Extras["prebuild"].(map[string]interface{})["dockerfile"] = base64.StdEncoding.EncodeToString(data)
}

func HandleStartLoggingAtBlock(cmd *cobra.Command, args []string, conf *Config) {
	if !cmd.Flags().Changed("start-logging-at-block") { //Don't bother if not specified
		return
	}

	startBlock, err := cmd.Flags().GetInt("start-logging-at-block")
	if err != nil {
		log.Trace("there was an error with the flag")
	} else {
		conf.Meta["startBlock"] = startBlock
	}
}

func HandleResources(cmd *cobra.Command, args []string, conf *Config) (givenCPU bool, givenMem bool) {
	givenCPU = cmd.Flags().Changed("cpus")
	givenMem = cmd.Flags().Changed("memory")

	cpus, err := cmd.Flags().GetString("cpus")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	memory, err := cmd.Flags().GetString("memory")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	if cpus != "0" {
		conf.Resources[0].Cpus = cpus
	}
	if memory != "0" {
		conf.Resources[0].Memory = memory
	}
	return
}

func HandleRepoBuild(cmd *cobra.Command, args []string, conf *Config) {
	if !cmd.Flags().Changed("git-repo") {
		return
	}
	if conf.Extras == nil {
		conf.Extras = map[string]interface{}{}
	}

	if _, ok := conf.Extras["prebuild"]; !ok {
		conf.Extras["prebuild"] = map[string]interface{}{}
	}
	conf.Extras["prebuild"].(map[string]interface{})["build"] = true

	repo, err := cmd.Flags().GetString("git-repo")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	conf.Extras["prebuild"].(map[string]interface{})["repo"] = repo
	if cmd.Flags().Changed("git-repo-branch") {
		branch, err := cmd.Flags().GetString("git-repo-branch")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		log.Trace("given a git repo branch")
		conf.Extras["prebuild"].(map[string]interface{})["branch"] = branch
	}
}

func HandlePortMapping(cmd *cobra.Command, args []string, conf *Config) {
	if !cmd.Flags().Changed("expose-port-mapping") {
		return
	}
	portMapping, err := cmd.Flags().GetStringSlice("expose-port-mapping")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	firstResources := conf.Resources[0]
	for conf.Nodes > len(conf.Resources) {
		conf.Resources = append(conf.Resources, firstResources)
	}
	parsedPortMapping, err := util.ParseIntToStringSlice(portMapping)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	for node, mappings := range parsedPortMapping {
		conf.Resources[node].Ports = mappings
		log.WithFields(log.Fields{"node": node, "ports": mappings}).Trace("adding the port mapping")
	}
}

/*func HandleExposeAllBuildFlag()
cmd.Flags().Int32Slice("expose-all",[],"expose a port linearly for all nodes")*/
