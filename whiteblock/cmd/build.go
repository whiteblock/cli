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
	serversFlag    string
	blockchainFlag string
	nodesFlag      int
	paramsFile     string
	validators     int
	optionsFlag    map[string]string
	envFlag        map[string]string
)

func buildAttach(buildId string) {
	buildListener(buildId)
	err := util.WriteStore(".previous_build_id", []byte(buildId))
	util.DeleteStore(".in_progress_build_id")

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

	fmt.Println("Build Started Successfully.")
	fmt.Printf("Testnet ID : %v\n", buildReply)

	//Store the in progress builds temporary id until the build finishes
	err = util.WriteStore(".in_progress_build_id", []byte(buildReply.(string)))
	if err != nil {
		util.PrintErrorFatal(err)
	}

	buildAttach(buildReply.(string))
}

func Build(cmd *cobra.Command, args []string, isAppend bool) {
	var err error
	util.CheckArguments(cmd, args, 0, 0)
	buildConf, _ := getPreviousBuild() //Errors are ok with this.
	blockchainEnabled := len(blockchainFlag) > 0
	nodesEnabled := nodesFlag > 0

	defaultCpus := ""
	defaultMemory := ""
	buildConf.Resources = []build.Resources{build.Resources{Cpus: "", Memory: ""}}
	buildConf.Params = map[string]interface{}{}
	buildConf.Extras = map[string]interface{}{}
	buildConf.Meta = map[string]interface{}{}

	previousYesAll, err := cmd.Flags().GetBool("yes")
	if err != nil {
		util.PrintErrorFatal(err)
	}

	cpusEnabled, memoryEnabled := build.HandleResources(cmd, args, &buildConf)

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
	if !cpusEnabled {
		allowEmpty = append(allowEmpty, true)
		buildOpt = append(buildOpt, "cpus"+tern((defaultCpus == ""), "(empty for no limit)", " ("+defaultCpus+")"))
		defOpt = append(defOpt, fmt.Sprintf(defaultCpus))
	}
	if !memoryEnabled {
		allowEmpty = append(allowEmpty, true)
		buildOpt = append(buildOpt, "memory"+tern((defaultMemory == ""), "(empty for no limit)", " ("+defaultMemory+")"))
		defOpt = append(defOpt, fmt.Sprintf(defaultMemory))
	}

	buildArr := []string{}
	if os.Stdin == nil && len(buildOpt) > 0 {
		fmt.Println("Would drop into build wizard but is a non interactive context")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)

	for i := 0; i < len(buildOpt); i++ {
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
			fmt.Println("Value required")
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

	if !cpusEnabled {
		buildConf.Resources[0].Cpus = buildArr[offset]
		offset++
	}
	if !memoryEnabled {
		buildConf.Resources[0].Memory = buildArr[offset]
		//offset++
	}

	if len(serversFlag) > 0 {
		serversInter := strings.Split(serversFlag, ",")
		buildConf.Servers = []int{}
		for _, serverStr := range serversInter {
			serverNum, err := strconv.Atoi(serverStr)
			if err != nil {
				util.InvalidInteger("servers", serverStr, true)
			}
			buildConf.Servers = append(buildConf.Servers, serverNum)
		}
	} else if len(buildConf.Servers) == 0 {
		buildConf.Servers = getServer()
	}

	options := <-optionsChannel //Currently has a negative impact but will be positive in the future
	if validators < 0 && hasParam(options, "validators") && !isAppend {
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
					util.PrintStringError(err.Error())
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
	build.HandlePullFlag(cmd, args, &buildConf)
	build.HandleForceUnlockFlag(cmd, args, &buildConf)
	build.HandleDockerAuthFlags(cmd, args, &buildConf)
	build.HandleSSHOptions(cmd, args, &buildConf)
	build.HandleDockerfile(cmd, args, &buildConf)
	build.HandleRepoBuild(cmd, args, &buildConf)
	if !isAppend {
		build.HandleStartLoggingAtBlock(cmd, args, &buildConf)
	}

	build.HandlePortMapping(cmd, args, &buildConf)
	log.WithFields(log.Fields{"build": buildConf}).Trace("sending the build request")
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
		buildId, err := util.ReadStore(".in_progress_build_id")
		if err != nil || len(buildId) == 0 {
			fmt.Println("No in progress build found. Use build command to deploy a blockchain.")
			os.Exit(1)
		}
		buildAttach(string(buildId))
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

		fmt.Println(util.Prettypi(prevBuild))
		if previousYesAll || util.YesNoPrompt("Build from previous?") {
			fmt.Println("building from previous configuration")
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
		buildId, err := util.ReadStore(".in_progress_build_id")
		if err != nil || len(buildId) == 0 {
			fmt.Println("No in-progress build found. Use build command to deploy a blockchain.")
			os.Exit(1)
		}
		defer util.DeleteStore(".in_progress_build_id")
		util.JsonRpcCallAndPrint("stop_build", []string{string(buildId)})
	},
}

var buildFreezeCmd = &cobra.Command{
	Use:     "freeze",
	Aliases: []string{"pause"},
	Short:   "Pause a build",
	Long:    "Pause a build",
	Run: func(cmd *cobra.Command, args []string) {
		buildId, err := util.ReadStore(".in_progress_build_id")
		if err != nil || len(buildId) == 0 {
			fmt.Println("No in-progress build found. Use build command to deploy a blockchain.")
			os.Exit(1)
		}
		util.JsonRpcCallAndPrint("freeze_build", []string{string(buildId)})
	},
}

var buildUnfreezeCmd = &cobra.Command{
	Use:     "unfreeze",
	Aliases: []string{"thaw", "resume"},
	Short:   "Unpause a build",
	Long:    "Unpause a build",
	Run: func(cmd *cobra.Command, args []string) {
		buildId, err := util.ReadStore(".in_progress_build_id")
		if err != nil || len(buildId) == 0 {
			fmt.Println("No in-progress build found. Use build command to deploy a blockchain.")
			os.Exit(1)
		}
		util.JsonRpcCallAndPrint("unfreeze_build", []string{string(buildId)})
		buildAttach(string(buildId))
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

func addBuildFlagsToCommand(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&serversFlag, "servers", "s", "", "display server options")
	cmd.Flags().BoolP("yes", "y", false, "Yes to all prompts. Evokes default parameters.")
	cmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	cmd.Flags().IntVarP(&nodesFlag, "nodes", "n", 0, "specify number of nodes")
	cmd.Flags().StringP("cpus", "c", "", "specify number of cpus")
	cmd.Flags().StringP("memory", "m", "", "specify memory allocated")
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
	cmd.Flags().UintSlice("expose-all", []uint{}, "expose a port linearly for all nodes")
	//META FLAGS
	cmd.Flags().Int("start-logging-at-block", 0, "specify a later block number to start at")
}

func init() {
	addBuildFlagsToCommand(buildCmd)
	addBuildFlagsToCommand(buildAppendCmd)

	previousCmd.Flags().BoolP("yes", "y", false, "Yes to all prompts. Evokes default parameters.")

	buildCmd.AddCommand(previousCmd, buildAppendCmd, buildStopCmd, buildAttachCmd, buildFreezeCmd, buildUnfreezeCmd)
	RootCmd.AddCommand(buildCmd)
}
