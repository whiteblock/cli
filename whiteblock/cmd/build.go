package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/cmd/build"
	"github.com/whiteblock/cli/whiteblock/util"
	"os"
	"strconv"
	"strings"
)

var (
	blockchainFlag string
	nodesFlag      int
	paramsFile     string
	validators     int
	optionsFlag    map[string]string
	envFlag        map[string]string
)

func buildAttach(buildID string) {
	buildListener(buildID)
	err := util.Set("previous_build_id", buildID)
	util.Delete("in_progress_build_id")

	if err != nil {
		util.PrintErrorFatal(err)
	}
}

func buildStart(buildConfig interface{}, isAppend bool) {
	var buildReply interface{}
	var err error
	if isAppend {
		buildReply, err = getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		_, err = util.JsonRpcCall("add_nodes", []interface{}{buildReply, buildConfig})
		if err != nil {
			util.PrintErrorFatal(err)
		}
	} else {
		buildReply, err = util.JsonRpcCall("build", buildConfig)
		if err != nil {
			util.PrintErrorFatal(err)
		}
	}

	util.Print("Build Started Successfully.")
	util.Printf("Testnet ID : %v\n", buildReply)

	//Store the in progress builds temporary id until the build finishes
	err = util.Set("in_progress_build_id", buildReply.(string))
	if err != nil {
		util.PrintErrorFatal(err)
	}

	buildAttach(buildReply.(string))
}

func Build(cmd *cobra.Command, args []string, isAppend bool) {
	var err error
	util.CheckArguments(cmd, args, 0, 0)
	buildConf, _ := getPreviousBuild() //Errors are ok with this.

	previousNumberNodes := 0
	if isAppend {
		nodes, err := GetNodes()
		log.WithFields(log.Fields{"nodes": buildConf.Nodes,
			"err": err}).Debug("getting node number from previous build")
		previousNumberNodes = len(nodes)
	}

	blockchainEnabled := len(blockchainFlag) > 0
	nodesEnabled := nodesFlag > 0

	buildConf.Resources = []build.Resources{build.Resources{Cpus: "", Memory: ""}}
	buildConf.Params = map[string]interface{}{}
	buildConf.Extras = map[string]interface{}{}
	buildConf.Meta = map[string]interface{}{}

	previousYesAll, err := cmd.Flags().GetBool("yes")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	build.HandleResources(cmd, args, &buildConf)

	buildOpt := []string{}
	defOpt := []string{}
	allowEmpty := []bool{}

	if !blockchainEnabled {
		allowEmpty = append(allowEmpty, false)
		buildOpt = append(buildOpt, "blockchain"+tern((len(buildConf.Blockchain) == 0), "", " ("+buildConf.Blockchain+")"))
		defOpt = append(defOpt, fmt.Sprintf(buildConf.Blockchain))
	}
	if !nodesEnabled {
		allowEmpty = append(allowEmpty, false)
		buildOpt = append(buildOpt, fmt.Sprintf("nodes(%d)", buildConf.Nodes))
		defOpt = append(defOpt, fmt.Sprintf("%d", buildConf.Nodes))
	}

	buildArr := []string{}
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < len(buildOpt); i++ {
		if !util.IsTTY() {
			util.PrintErrorFatal("missing build parameters and couldn't prompt")
		}
		fmt.Print(buildOpt[i] + ": ")
		if !scanner.Scan() {
			util.PrintErrorFatal(scanner.Err())
		}

		text := scanner.Text()
		if len(text) != 0 {
			buildArr = append(buildArr, text)
		} else if len(defOpt[i]) != 0 || allowEmpty[i] {
			buildArr = append(buildArr, defOpt[i])
		} else {
			i--
			util.Print("Value required")
			continue
		}
	}

	offset := 0
	if blockchainEnabled {
		buildConf.Blockchain = strings.ToLower(blockchainFlag)
	} else {
		buildConf.Blockchain = strings.ToLower(buildArr[offset])
		offset++
	} //Final blockchain definition. Will need to start another round of prompting
	optionsChannel := make(chan [][]string, 1)
	go func() {
		if buildConf.Blockchain == "generic" {
			optionsChannel <- [][]string{}
			return
		}
		opt, err := fetchParams(buildConf.Blockchain)
		if err != nil {
			util.PrintErrorFatal(err)
		}
		optionsChannel <- opt
	}()

	if nodesEnabled {
		buildConf.Nodes = nodesFlag
	} else {
		buildConf.Nodes, err = strconv.Atoi(buildArr[offset])
		if err != nil {
			util.InvalidInteger("nodes", buildArr[offset], true)
		}
		offset++
	}

	options := <-optionsChannel //Currently has a negative impact but will be positive in the future
	if validators < 0 && hasParam(options, "validators") && !isAppend {
		if !util.IsTTY() {
			util.PrintErrorFatal("missing validators and couldn't prompt")
		}
		fmt.Print("validators: ")
		if !scanner.Scan() {
			util.PrintErrorFatal(scanner.Err())
		}
		text := scanner.Text()
		validators, err = strconv.Atoi(text)
		if err != nil {
			util.PrintErrorFatal(err)
		}
	}
	build.HandleImageFlag(cmd, args, &buildConf)
	if optionsFlag != nil {
		buildConf.Params, err = processOptions(optionsFlag, options)
		if err != nil {
			util.PrintErrorFatal(err)
		}
	} else if len(paramsFile) != 0 {
		f, err := os.Open(paramsFile)
		if err != nil {
			util.PrintErrorFatal(err)
		}

		decoder := json.NewDecoder(f)
		decoder.UseNumber()
		err = decoder.Decode(&buildConf.Params)
		if err != nil {
			util.PrintErrorFatal(err)
		}
	} else if !previousYesAll && !util.YesNoPrompt("Use default parameters?") {
		if !util.IsTTY() {
			util.PrintErrorFatal("not a tty")
		}
		//PARAMS
		for i := 0; i < len(options); i++ {
			key := options[i][0]
			key_type := options[i][1]

			fmt.Printf("%s (%s): ", key, key_type)
			scanner.Scan()
			text := scanner.Text()
			if len(text) == 0 {
				continue
			}
			switch key_type {
			case "string":
				//needs to have filtering
				buildConf.Params[key] = text
			case "[]string":
				preprocessed := strings.Replace(text, " ", ",", -1)
				buildConf.Params[key] = strings.Split(preprocessed, ",")
			case "int":
				val, err := strconv.ParseInt(text, 0, 64)
				if err != nil {
					util.InvalidInteger(key, text, false)
					i--
					continue
				}
				buildConf.Params[key] = val
			case "bool":
				val, err := util.GetAsBool(text)
				if err != nil {
					util.PrintErrorFatal(err)
					i--
					continue
				}
				buildConf.Params[key] = val
			}
		}
	}
	if validators >= 0 {
		buildConf.Params["validators"] = validators
	}
	build.HandleFilesFlag(cmd, args, &buildConf)

	if envFlag != nil {
		buildConf.Environments, err = processEnv(envFlag, buildConf.Nodes)
		if err != nil {
			util.PrintErrorFatal(err)
		}
	}

	fbg, err := cmd.Flags().GetBool("freeze-before-genesis")
	if err == nil && fbg {
		buildConf.Extras["freezeAfterInfrastructure"] = true
	}

	build.HandleServersFlag(cmd, args, &buildConf)
	build.HandlePullFlag(cmd, args, &buildConf)
	build.HandleForceUnlockFlag(cmd, args, &buildConf)
	build.HandleDockerAuthFlags(cmd, args, &buildConf)
	build.HandleSSHOptions(cmd, args, &buildConf)
	build.HandleDockerfile(cmd, args, &buildConf)
	build.HandleRepoBuild(cmd, args, &buildConf)
	build.HandleBoundCPUs(cmd, args, &buildConf)
	if !isAppend {
		build.HandleStartLoggingAtBlock(cmd, args, &buildConf)
	}

	build.HandlePortMapping(cmd, args, &buildConf)
	build.HandleExposeAllBuildFlag(cmd, args, &buildConf, previousNumberNodes)

	log.WithFields(log.Fields{"build": buildConf, "dest": conf.ServerAddr, "api": conf.APIURL}).Trace("sending the build request")
	build.SanitizeBuild(&buildConf)
	buildStart(buildConf, isAppend)
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create", "buidl"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: "Build will create and deploy a blockchain and the specified number of nodes." +
		" Each node will be instantiated in its own container and will interact" +
		" individually as a participant of the specified network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		Build(cmd, args, false)
	},
}

var buildAttachCmd = &cobra.Command{
	Use:     "attach",
	Aliases: []string{"resume"},
	Short:   "Build a blockchain using previous configurations",
	Long:    "\nAttach to a current in progress build process\n",

	Run: func(cmd *cobra.Command, args []string) {
		var buildID string
		err := util.GetP("in_progress_build_id", &buildID)
		if err != nil || len(buildID) == 0 {
			util.PrintErrorFatal("No in progress build found. Use build command to deploy a blockchain.")
		}
		buildAttach(buildID)
	},
}

var previousCmd = &cobra.Command{
	Use:     "previous",
	Aliases: []string{"prev"},
	Short:   "Build a blockchain using previous configurations",
	Long:    "\nBuild previous will recreate and deploy the previously built blockchain and specified number of nodes.\n",

	Run: func(cmd *cobra.Command, args []string) {

		prevBuild, err := getPreviousBuild()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		previousYesAll, err := cmd.Flags().GetBool("yes")
		if err != nil {
			util.PrintErrorFatal(err)
		}
		util.Print(prevBuild)
		if previousYesAll || util.YesNoPrompt("Build from previous?") {
			util.Print("building from previous configuration")
			buildStart(prevBuild, false)
			return
		}
	},
}

var buildStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"cancel"},
	Short:   "Stops the current build",
	Long:    "\nBuild stops the current building process.\n",

	Run: func(cmd *cobra.Command, args []string) {
		var buildID string
		err := util.GetP("in_progress_build_id", &buildID)
		if err != nil || len(buildID) == 0 {
			util.PrintErrorFatal("No in-progress build found. Use build command to deploy a blockchain.")
		}
		defer util.Delete("in_progress_build_id")
		util.JsonRpcCallAndPrint("stop_build", []interface{}{buildID})
	},
}

var buildFreezeCmd = &cobra.Command{
	Use:     "freeze",
	Aliases: []string{"pause"},
	Short:   "Pause a build",
	Long:    "Pause a build",
	Run: func(cmd *cobra.Command, args []string) {
		var buildID string
		err := util.GetP("in_progress_build_id", &buildID)
		if err != nil || len(buildID) == 0 {
			util.PrintErrorFatal("No in-progress build found. Use build command to deploy a blockchain.")
		}
		util.JsonRpcCallAndPrint("freeze_build", []string{buildID})
	},
}

var buildUnfreezeCmd = &cobra.Command{
	Use:     "unfreeze",
	Aliases: []string{"thaw", "resume"},
	Short:   "Unpause a build",
	Long:    "Unpause a build",
	Run: func(cmd *cobra.Command, args []string) {
		var buildID string
		err := util.GetP("in_progress_build_id", &buildID)
		if err != nil || len(buildID) == 0 {
			util.PrintErrorFatal("No in-progress build found. Use build command to deploy a blockchain.")
		}
		util.JsonRpcCallAndPrint("unfreeze_build", []string{buildID})
		buildAttach(buildID)
	},
}

var buildAppendCmd = &cobra.Command{
	Use:   "append",
	Short: "Build a blockchain using image and deploy nodes",
	Long: "Build will create and deploy a blockchain and the specified number of nodes." +
		" Each node will be instantiated in its own container and will interact" +
		" individually as a participant of the specified network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		Build(cmd, args, true)
	},
}

func addBuildFlagsToCommand(cmd *cobra.Command, isAppend bool) {
	cmd.Flags().IntSliceP("servers", "s", []int{}, "manually choose the server options")
	cmd.Flags().BoolP("yes", "y", false, "Yes to all prompts. Evokes default parameters.")
	cmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	cmd.Flags().IntVarP(&nodesFlag, "nodes", "n", 0, "specify number of nodes")
	cmd.Flags().StringP("cpus", "c", "0", "specify number of cpus")
	cmd.Flags().StringP("memory", "m", "0", "specify memory allocated")
	cmd.Flags().StringVarP(&paramsFile, "file", "f", "", "parameters file")
	cmd.Flags().IntVarP(&validators, "validators", "v", -1, "set the number of validators")
	cmd.Flags().StringSliceP("image", "i", []string{}, "image tag")
	cmd.Flags().StringToStringVarP(&optionsFlag, "option", "o", nil, "blockchain specific options")
	cmd.Flags().StringToStringVarP(&envFlag, "env", "e", nil, "set environment variables for the nodes")
	cmd.Flags().StringSliceP("template", "t", nil, "set a custom file template")

	cmd.Flags().String("docker-username", "", "docker auth username")
	cmd.Flags().String("docker-password", "", "docker auth password. Note: this will be stored unencrypted while the build is in progress")
	cmd.Flags().StringSlice("user-ssh-key", []string{}, "add an additional ssh key as authorized for the nodes."+
		" Takes a file containing an ssh public key")

	cmd.Flags().Bool("force-docker-pull", false, "Manually pull the image before the build")
	cmd.Flags().Bool("force-unlock", false, "Forcefully stop and unlock the build process")
	cmd.Flags().Bool("freeze-before-genesis", false, "indicate that the build should freeze before starting the genesis ceremony")
	cmd.Flags().String("dockerfile", "", "build from a dockerfile")
	cmd.Flags().StringSliceP("expose-port-mapping", "p", nil, "expose a port to the outside world -p 0=8545:8546")

	cmd.Flags().String("git-repo", "", "build from a git repo")
	cmd.Flags().String("git-repo-branch", "", "specify the branch to build from in a git repo")
	cmd.Flags().IntSlice("expose-all", []int{}, "expose a port linearly for all nodes")
	//META FLAGS
	if !isAppend {
		cmd.Flags().Int("start-logging-at-block", 0, "specify a later block number to start at")
		cmd.Flags().Int("bound-cpus", -1, "specify number of bound cpus")
	}

}

func init() {
	addBuildFlagsToCommand(buildCmd, false)
	addBuildFlagsToCommand(buildAppendCmd, true)

	previousCmd.Flags().BoolP("yes", "y", false, "Yes to all prompts. Evokes default parameters.")

	buildCmd.AddCommand(previousCmd, buildAppendCmd, buildStopCmd, buildAttachCmd, buildFreezeCmd, buildUnfreezeCmd)
	RootCmd.AddCommand(buildCmd)
}
