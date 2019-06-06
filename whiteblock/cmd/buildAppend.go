package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	util "github.com/whiteblock/cli/whiteblock/util"
	"os"
	"strings"
)

var buildAppendCmd = &cobra.Command{
	Use: "append",
	//Aliases: []string{"init", "create", "buidl"},
	Short: "Build a blockchain using image and deploy nodes",
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

		if blockchainEnabled {
			buildConf.Blockchain = strings.ToLower(blockchainFlag)
		}

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
			buildConf.Nodes = 1
		}
		options := <-optionsChannel //Currently has a negative impact but will be positive in the future

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
		}
		handleFilesFlag(cmd, args, &buildConf)

		if envFlag != nil {
			buildConf.Environments, err = processEnv(envFlag, buildConf.Nodes)
			if err != nil {
				util.PrintErrorFatal(err)
			}
		}
		if validators >= 0 {
			buildConf.Params["validators"] = validators
		}
		fbg, err := cmd.Flags().GetBool("freeze-before-genesis")
		if err == nil && fbg {
			buildConf.Extras["freezeAfterInfrastructure"] = true
		}
		handlePullFlag(cmd, args, &buildConf)
		handleDockerAuthFlags(cmd, args, &buildConf)
		//fmt.Printf("%+v\n", buildConf)
		testnetID, err := getPreviousBuildId()
		if err != nil {
			util.PrintErrorFatal(err)
		}
		//fmt.Printf("%+v\n", buildConf)
		_, err = jsonRpcCall("add_nodes", []interface{}{testnetID, buildConf})
		if err != nil {
			util.PrintErrorFatal(err)
		}
		fmt.Printf("Adding Nodes Started successfully: %v\n", testnetID)

		//Store the in progress builds temporary id until the build finishes
		err = util.WriteStore(".in_progress_build_id", []byte(testnetID))
		if err != nil {
			util.PrintErrorFatal(err)
		}

		buildAttach(testnetID)

		removeSmartContracts()
	},
}

func init() {

	buildAppendCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	buildAppendCmd.Flags().IntVarP(&nodesFlag, "nodes", "n", 0, "specify number of nodes")
	buildAppendCmd.Flags().StringVarP(&cpusFlag, "cpus", "c", "", "specify number of cpus")
	buildAppendCmd.Flags().StringVarP(&memoryFlag, "memory", "m", "", "specify memory allocated")
	buildAppendCmd.Flags().StringVarP(&paramsFile, "file", "f", "", "parameters file")
	buildAppendCmd.Flags().IntVarP(&validators, "validators", "v", -1, "set the number of validators")
	buildAppendCmd.Flags().StringSliceP("image", "i", []string{}, "image tag")
	buildAppendCmd.Flags().StringToStringVarP(&optionsFlag, "option", "o", nil, "blockchain specific options")
	buildAppendCmd.Flags().StringToStringVarP(&envFlag, "env", "e", nil, "set environment variables for the nodes")
	buildAppendCmd.Flags().StringSliceP("template", "t", nil, "set a custom file template")

	buildAppendCmd.Flags().String("docker-username", "", "docker auth username")
	buildAppendCmd.Flags().String("docker-password", "", "docker auth password. Note: this will be stored unencrypted while the build is in progress")
	buildAppendCmd.Flags().Bool("force-docker-pull", false, "Manually pull the image before the build")
	buildAppendCmd.Flags().Bool("freeze-before-genesis", false, "indicate that the build should freeze before starting the genesis ceremony")
	buildCmd.AddCommand(buildAppendCmd)
}
