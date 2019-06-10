package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	util "github.com/whiteblock/cli/whiteblock/util"
	"os"
	"strconv"
	"strings"
)

var (
	previousYesAll bool
	serversFlag    string
	blockchainFlag string
	nodesFlag      int
	cpusFlag       string
	memoryFlag     string
	paramsFile     string
	validators     int
	optionsFlag    map[string]string
	envFlag        map[string]string
)

type Config struct {
	Servers      []int                  `json:"servers"`
	Blockchain   string                 `json:"blockchain"`
	Nodes        int                    `json:"nodes"`
	Images       []string               `json:"images"`
	Resources    []Resources            `json:"resources"`
	Params       map[string]interface{} `json:"params"`
	Environments []map[string]string    `json:"environments"`
	Files        []map[string]string    `json:"files"`
	Logs         []map[string]string    `json:"logs"`
	Extras       map[string]interface{} `json:"extras"`
}

type Resources struct {
	Cpus   string `json:"cpus"`
	Memory string `json:"memory"`
}

func buildAttach(buildId string) {
	buildListener(buildId)
	err := util.WriteStore(".previous_build_id", []byte(buildId))
	util.DeleteStore(".in_progress_build_id")

	if err != nil {
		util.PrintErrorFatal(err)
	}
}

func build(buildConfig interface{}) {
	buildReply, err := jsonRpcCall("build", buildConfig)
	if err != nil {
		util.PrintErrorFatal(err)
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

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create", "buidl"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: "Build will create and deploy a blockchain and the specified number of nodes." +
		" Each node will be instantiated in its own container and will interact" +
		" individually as a participant of the specified network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		var err error
		util.CheckArguments(cmd, args, 0, 0)
		buildConf, _ := getPreviousBuild() //Errors are ok with this.
		blockchainEnabled := len(blockchainFlag) > 0
		nodesEnabled := nodesFlag > 0
		cpusEnabled := len(cpusFlag) != 0
		memoryEnabled := len(memoryFlag) != 0

		defaultCpus := ""
		defaultMemory := ""

		if buildConf.Resources != nil && len(buildConf.Resources) > 0 {
			defaultCpus = string(buildConf.Resources[0].Cpus)
			defaultMemory = "" //string(buildConf.Resources[0].Memory)
		} else if buildConf.Resources == nil {
			buildConf.Resources = []Resources{Resources{}}
		}

		buildConf.Params = map[string]interface{}{}
		buildConf.Extras = map[string]interface{}{}

		if cpusFlag == "0" {
			cpusFlag = ""
		} else if cpusEnabled {
			buildConf.Resources[0].Cpus = cpusFlag
		}

		if memoryFlag == "0" {
			memoryFlag = ""
		} else if memoryEnabled {
			buildConf.Resources[0].Memory = memoryFlag
		}

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
		if validators < 0 && hasParam(options, "validators") {
			fmt.Print("validators: ")
			scanner.Scan()
			text := scanner.Text()
			validators, err = strconv.Atoi(text)
			if err != nil {
				util.PrintErrorFatal(err)
			}
		}
		handleImageFlag(cmd, args, &buildConf)
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

			//scanner := bufio.NewScanner(os.Stdin)

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
				}
			}
		}
		if validators >= 0 {
			buildConf.Params["validators"] = validators
		}
		handleFilesFlag(cmd, args, &buildConf)

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
		handlePullFlag(cmd, args, &buildConf)
		handleForceUnlockFlag(cmd, args, &buildConf)
		handleDockerAuthFlags(cmd, args, &buildConf)
		handleSSHOptions(cmd, args, &buildConf)
		//fmt.Printf("%+v\n", buildConf)
		build(buildConf)
		removeSmartContracts()
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
		fmt.Println(prettypi(prevBuild))
		if previousYesAll || util.YesNoPrompt("Build from previous?") {
			fmt.Println("building from previous configuration")
			build(prevBuild)
			removeSmartContracts()
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
		jsonRpcCallAndPrint("stop_build", []string{string(buildId)})
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
		jsonRpcCallAndPrint("freeze_build", []string{string(buildId)})
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
		jsonRpcCallAndPrint("unfreeze_build", []string{string(buildId)})
		buildAttach(string(buildId))
	},
}

func init() {
	buildCmd.Flags().StringVarP(&serversFlag, "servers", "s", "", "display server options")
	buildCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts. Evokes default parameters.")
	buildCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	buildCmd.Flags().IntVarP(&nodesFlag, "nodes", "n", 0, "specify number of nodes")
	buildCmd.Flags().StringVarP(&cpusFlag, "cpus", "c", "", "specify number of cpus")
	buildCmd.Flags().StringVarP(&memoryFlag, "memory", "m", "", "specify memory allocated")
	buildCmd.Flags().StringVarP(&paramsFile, "file", "f", "", "parameters file")
	buildCmd.Flags().IntVarP(&validators, "validators", "v", -1, "set the number of validators")
	buildCmd.Flags().StringSliceP("image", "i", []string{}, "image tag")
	buildCmd.Flags().StringToStringVarP(&optionsFlag, "option", "o", nil, "blockchain specific options")
	buildCmd.Flags().StringToStringVarP(&envFlag, "env", "e", nil, "set environment variables for the nodes")
	buildCmd.Flags().StringSliceP("template", "t", nil, "set a custom file template")

	buildCmd.Flags().String("docker-username", "", "docker auth username")
	buildCmd.Flags().String("docker-password", "", "docker auth password. Note: this will be stored unencrypted while the build is in progress")
	buildCmd.Flags().StringSlice("user-ssh-key", []string{}, "add an additional ssh key as authorized for the nodes."+
		" Takes a file containing an ssh public key")

	buildCmd.Flags().Bool("force-docker-pull", false, "Manually pull the image before the build")
	buildCmd.Flags().Bool("force-unlock", false, "Forcefully stop and unlock the build process")
	buildCmd.Flags().Bool("freeze-before-genesis", false, "indicate that the build should freeze before starting the genesis ceremony")

	previousCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts. Evokes default parameters.")

	buildCmd.AddCommand(previousCmd, buildStopCmd, buildAttachCmd, buildFreezeCmd, buildUnfreezeCmd)
	RootCmd.AddCommand(buildCmd)
}
