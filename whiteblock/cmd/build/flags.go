package build

import (
	"encoding/base64"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"strconv"
	"strings"
)

func HandleForceUnlockFlag(cmd *cobra.Command, args []string, bconf *Config) {

	fbg, err := cmd.Flags().GetBool("force-unlock")
	if err == nil && fbg {
		bconf.Extras["forceUnlock"] = true
	}
}

func HandlePullFlag(cmd *cobra.Command, args []string, bconf *Config) {
	_, ok := bconf.Extras["prebuild"]
	if !ok {
		bconf.Extras["prebuild"] = map[string]interface{}{}
	}
	fbg, err := cmd.Flags().GetBool("force-docker-pull")
	if err == nil && fbg {
		bconf.Extras["prebuild"].(map[string]interface{})["pull"] = true
	}
}

func HandleDockerAuthFlags(cmd *cobra.Command, args []string, bconf *Config) {
	if cmd.Flags().Changed("docker-password") != cmd.Flags().Changed("docker-username") {
		if cmd.Flags().Changed("docker-password") {
			util.PrintErrorFatal("you must also provide --docker-password with --docker-username")
		}
		util.PrintErrorFatal("you must also provide --docker-username with --docker-password")
	}
	if !cmd.Flags().Changed("docker-password") {
		return //The auth flags have not been set
	}

	_, ok := bconf.Extras["prebuild"]
	if !ok {
		bconf.Extras["prebuild"] = map[string]interface{}{}
	}
	username, err := cmd.Flags().GetString("docker-username")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	password, err := cmd.Flags().GetString("docker-password")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	bconf.Extras["prebuild"].(map[string]interface{})["auth"] = map[string]string{
		"username": username,
		"password": password,
	}

}

func HandleImageFlag(cmd *cobra.Command, args []string, bconf *Config) {

	imageFlag, err := cmd.Flags().GetStringSlice("image")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	bconf.Images = make([]string, bconf.Nodes)
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
	baseImage := getImage(bconf.Blockchain, "stable", imgDefault)

	for i := 0; i < bconf.Nodes; i++ {

		bconf.Images[i] = baseImage
		image, exists := images[i]
		if exists {
			log.WithFields(log.Fields{"image": image}).Trace("image exists")
			bconf.Images[i] = image
		}
	}
}

func HandleFilesFlag(cmd *cobra.Command, args []string, bconf *Config) {
	filesFlag, err := cmd.Flags().GetStringSlice("template")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	if filesFlag == nil {
		return
	}

	bconf.Files = make([]map[string]string, bconf.Nodes)
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
		if index < 0 || index >= bconf.Nodes {
			util.PrintErrorFatal(fmt.Errorf("Index is out of range for -t flag"))
		}
		bconf.Files[index] = map[string]string{}
		bconf.Files[index][tuple[1]] = base64.StdEncoding.EncodeToString(data)
	}

	if bconf.Extras == nil {
		bconf.Extras = map[string]interface{}{}
	}
	if _, ok := bconf.Extras["defaults"]; !ok {
		bconf.Extras["defaults"] = map[string]interface{}{}
	}
	bconf.Extras["defaults"].(map[string]interface{})["files"] = defaults
}

func HandleSSHOptions(cmd *cobra.Command, args []string, bconf *Config) {
	if !cmd.Flags().Changed("user-ssh-key") { //Don't bother if not specified
		return
	}

	sshPubKeys, err := cmd.Flags().GetStringSlice("user-ssh-key")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	if bconf.Extras == nil {
		bconf.Extras = map[string]interface{}{}
	}
	if _, ok := bconf.Extras["postbuild"]; !ok {
		bconf.Extras["postbuild"] = map[string]interface{}{}
	}
	if _, ok := bconf.Extras["postbuild"].(map[string]interface{})["ssh"]; !ok {
		bconf.Extras["postbuild"].(map[string]interface{})["ssh"] = map[string]interface{}{}
	}
	pubKeys := []string{}
	for _, pubKeyFile := range sshPubKeys {
		data, err := ioutil.ReadFile(pubKeyFile)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		pubKeys = append(pubKeys, string(data))
	}

	bconf.Extras["postbuild"].(map[string]interface{})["ssh"].(map[string]interface{})["pubKeys"] = pubKeys
}

func HandleDockerfile(cmd *cobra.Command, args []string, bconf *Config) {
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

	if bconf.Extras == nil {
		bconf.Extras = map[string]interface{}{}
	}

	if _, ok := bconf.Extras["prebuild"]; !ok {
		bconf.Extras["prebuild"] = map[string]interface{}{}
	}
	bconf.Extras["prebuild"].(map[string]interface{})["build"] = true
	bconf.Extras["prebuild"].(map[string]interface{})["dockerfile"] = base64.StdEncoding.EncodeToString(data)
}

func HandleStartLoggingAtBlock(cmd *cobra.Command, args []string, bconf *Config) {
	if !cmd.Flags().Changed("start-logging-at-block") { //Don't bother if not specified
		return
	}

	startBlock, err := cmd.Flags().GetInt("start-logging-at-block")
	if err != nil {
		log.Trace("there was an error with the flag")
	} else {
		bconf.Meta["startBlock"] = startBlock
	}
}

func HandleResources(cmd *cobra.Command, args []string, bconf *Config) (givenCPU bool, givenMem bool) {
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
		bconf.Resources[0].Cpus = cpus
	}
	if memory != "0" {
		bconf.Resources[0].Memory = memory
	}
	return
}

func HandleRepoBuild(cmd *cobra.Command, args []string, bconf *Config) {
	if !cmd.Flags().Changed("git-repo") {
		return
	}
	if bconf.Extras == nil {
		bconf.Extras = map[string]interface{}{}
	}

	if _, ok := bconf.Extras["prebuild"]; !ok {
		bconf.Extras["prebuild"] = map[string]interface{}{}
	}
	bconf.Extras["prebuild"].(map[string]interface{})["build"] = true

	repo, err := cmd.Flags().GetString("git-repo")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	bconf.Extras["prebuild"].(map[string]interface{})["repo"] = repo
	if cmd.Flags().Changed("git-repo-branch") {
		branch, err := cmd.Flags().GetString("git-repo-branch")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		log.Trace("given a git repo branch")
		bconf.Extras["prebuild"].(map[string]interface{})["branch"] = branch
	}
}

func addPortMapping(portMapping map[int][]string, bconf *Config) {
	firstResources := bconf.Resources[0]
	for bconf.Nodes > len(bconf.Resources) {
		bconf.Resources = append(bconf.Resources, firstResources)
	}
	for node, mappings := range portMapping {
		bconf.Resources[node].Ports = mappings
		log.WithFields(log.Fields{"node": node, "ports": mappings}).Trace("adding the port mapping")
	}
}

func HandlePortMapping(cmd *cobra.Command, args []string, bconf *Config) {
	if !cmd.Flags().Changed("expose-port-mapping") {
		return
	}
	portMapping, err := cmd.Flags().GetStringSlice("expose-port-mapping")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	parsedPortMapping, err := util.ParseIntToStringSlice(portMapping)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	addPortMapping(parsedPortMapping, bconf)

}

func HandleExposeAllBuildFlag(cmd *cobra.Command, args []string, bconf *Config, offset int) {
	if !cmd.Flags().Changed("expose-all") {
		return
	}
	portsToExpose, err := cmd.Flags().GetIntSlice("expose-all")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	portMapping := map[int][]string{}
	usedPort := map[int]bool{}
	for i := 0; i < bconf.Nodes; i++ {
		portMapping[i] = []string{}
		for _, portToExpose := range portsToExpose {
			portToBind := portToExpose + i + offset
			_, used := usedPort[portToBind]
			if used {
				util.PrintErrorFatal(
					fmt.Sprintf("would duplicate exposed port %d. Too many nodes to run auto expose", portToExpose))
			}
			portMapping[i] = append(portMapping[i], fmt.Sprintf("%d:%d", portToBind, portToExpose))
		}
	}
	addPortMapping(portMapping, bconf)
}

func HandleServersFlag(cmd *cobra.Command, args []string, bconf *Config) {
	if !cmd.Flags().Changed("servers") {
		bconf.Servers = getServer()
		return
	}
	servers, err := cmd.Flags().GetIntSlice("servers")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	bconf.Servers = servers
}

func HandleBoundCPUs(cmd *cobra.Command, args []string, bconf *Config) {
	if !cmd.Flags().Changed("bound-cpus") {
		return
	}
	firstResources := bconf.Resources[0]
	for bconf.Nodes > len(bconf.Resources) {
		bconf.Resources = append(bconf.Resources, firstResources)
	}
	numCPUs, err := cmd.Flags().GetInt("bound-cpus")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	cpuNo := 0
	for i := range bconf.Resources {
		bconf.Resources[i].BoundCPUs = []int{}
		for j := 0; j < numCPUs; j++ {
			bconf.Resources[i].BoundCPUs = append(bconf.Resources[i].BoundCPUs, cpuNo)
			cpuNo++
		}
	}
}
