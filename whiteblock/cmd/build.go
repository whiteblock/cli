package cmd

import (
	"bufio"
	"encoding/json"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"errors"
	"github.com/spf13/cobra"
	util "../util"
)

var (
	previousYesAll 	bool
	serversFlag    	string
	blockchainFlag 	string
	nodesFlag      	int
	cpusFlag       	string
	memoryFlag     	string
	paramsFile     	string
	validators     	int
	imageFlag      	string
	optionsFlag		map[string]string
	envFlag			map[string]string
	filesFlag		map[string]string
)

type Config struct {
	Servers    		[]int                  	`json:"servers"`
	Blockchain 		string                 	`json:"blockchain"`
	Nodes      		int                    	`json:"nodes"`
	Image      		string                 	`json:"image"`
	Resources  		[]Resources            	`json:"resources"`
	Params     		map[string]interface{} 	`json:"params"`
	Environments    []map[string]string    	`json:"environments"`
	Files			map[string]string		`json:"files"`
	
}

type Resources struct {
	Cpus   string 	`json:"cpus"`
	Memory string 	`json:"memory"`
}

func getPreviousBuildId() (string,error) {
	buildId,err := util.ReadStore(".previous_build_id")
	if err != nil || len(buildId) == 0 {
		return "",errors.New("No previous build. Use build command to deploy a blockchain.")
	}
	return string(buildId),nil
}

func getPreviousBuild() (Config,error) {
	buildId,err := getPreviousBuildId()
	if err != nil {
		return Config{},err
	}

	prevBuild, err := jsonRpcCall("get_build", []string{buildId})
	if err != nil {
		return Config{},err
	}

	tmp,err := json.Marshal(prevBuild)
	if err != nil {
		return Config{},err
	}

	var out Config
	err = json.Unmarshal(tmp,&out)
	return out,err
}

func buildAttach(buildId string) {
	buildListener(buildId)
	err := util.WriteStore(".previous_build_id",[]byte(buildId))
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
	fmt.Printf("Build Started successfully: %v\n", buildReply)

	//Store the in progress builds temporary id until the build finishes
	err = util.WriteStore(".in_progress_build_id",[]byte(buildReply.(string)))
	if err != nil {
		util.PrintErrorFatal(err)
	}
	
	buildAttach(buildReply.(string))
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
		idList = append(idList,serverID)
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

func getImage(blockchain string,imageType string,defaultImage string) string {
	cwd := os.Getenv("HOME")
	b, err := ioutil.ReadFile("/etc/whiteblock.json")
	if err != nil {
		b, err = ioutil.ReadFile(cwd + "/cli/etc/whiteblock.json")
		if err != nil {
			util.PrintErrorFatal(err)
		}
	}
	var cont map[string]map[string]map[string]map[string]string
	err = json.Unmarshal(b, &cont)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	// fmt.Println(cont["blockchains"][blockchain]["images"][image])
	if len(cont["blockchains"][blockchain]["images"][imageType]) != 0 {
		return cont["blockchains"][blockchain]["images"][imageType]
	} else if(len(defaultImage) > 0){
		return defaultImage
	}else{
		return "gcr.io/whiteblock/"+blockchain+":master"
	}
}

func removeSmartContracts() {
	cwd := os.Getenv("HOME")
	err := os.RemoveAll(cwd + "/smart-contracts/whiteblock/contracts.json")
	if err != nil {
		util.PrintErrorFatal(err)
	}
}

func processOptions(givenOptions map[string]string,format []interface{}) (map[string]interface{},error) {
	out := map[string]interface{}{}
	
	for _,kv := range format {
		_kv := kv.([]interface{})
		name := _kv[0].(string)
		key_type := _kv[1].(string)

		val,ok := givenOptions[name]
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
					return nil,err
				}
				out[name] = val
		}	
	}
	return out,nil
}
//-1 means for all
func processEnvKey(in string) (int,string) {
	node := -1
	index := 0
	for i,char := range in {
		if char < '0' || char > '9' {
			index = i
			break
		}
	}
	if index == 0 {
		return node,in
	}

	if index == len(in) {
		util.PrintStringError("Cannot have a numerical environment variable")
		os.Exit(1)
	}

	var err error
	node,err = strconv.Atoi(in[0:index])
	if err != nil {
		util.PrintErrorFatal(err)
	}
	return node,in[index:len(in)]
}

func processEnv(envVars map[string]string,nodes int) ([]map[string]string,error) {
	out := make([]map[string]string,nodes)
	for i,_ := range out {
		out[i] = make(map[string]string)
	}
	for k,v := range envVars {
		node,key := processEnvKey(k)
		if node == -1 {
			for i,_ := range out {
				out[i][key] = v
			}
			continue
		}
		out[node][key] = v
	}
	return out,nil
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: "Build will create and deploy a blockchain and the specified number of nodes." +
		" Each node will be instantiated in its own container and will interact" +
		" individually as a participant of the specified network.\n",

	Run: func(cmd *cobra.Command, args []string) {

		
		util.CheckArguments(cmd,args,0,0)
		buildConf,err := getPreviousBuild()
		if err != nil {
			//util.PrintError(err)
		}
		blockchainEnabled := len(blockchainFlag) > 0
		nodesEnabled := nodesFlag > 0
		cpusEnabled := len(cpusFlag) != 0
		memoryEnabled := len(memoryFlag) != 0


		
		
		defaultCpus := ""
		defaultMemory := ""

		if buildConf.Resources != nil && len(buildConf.Resources) > 0 {
			defaultCpus = string(buildConf.Resources[0].Cpus)
			defaultMemory = string(buildConf.Resources[0].Memory)
		}else if buildConf.Resources == nil {
			buildConf.Resources = []Resources{Resources{}}
		}

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
			buildOpt = append(buildOpt, fmt.Sprintf("nodes(%d)",buildConf.Nodes))
			defOpt = append(defOpt, fmt.Sprintf("%d",buildConf.Nodes))
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
		scanner := bufio.NewScanner(os.Stdin)

		for i := 0; i < len(buildOpt); i++ {
			fmt.Print(buildOpt[i] + ": ")
			scanner.Scan()

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

		if len(serversFlag) > 0 {
			serversInter := strings.Split(serversFlag,",")
			buildConf.Servers = []int{}
			for _,serverStr := range serversInter {
				serverNum, err := strconv.Atoi(serverStr)
				if err != nil {
					util.InvalidInteger("servers",serverStr,true)
				}
				buildConf.Servers = append(buildConf.Servers,serverNum)
			}
		} else if len(buildConf.Servers) == 0 {
			buildConf.Servers = getServer()
		}

		offset := 0

		if blockchainEnabled {
			buildConf.Blockchain = blockchainFlag
		} else {
			buildConf.Blockchain = buildArr[offset]
			offset++
		}

		if nodesEnabled {
			buildConf.Nodes = nodesFlag
		} else {
			buildConf.Nodes, err = strconv.Atoi(buildArr[offset])
			if err != nil {
				util.InvalidInteger("nodes", buildArr[offset], true)
			}
			offset++
		}

		buildConf.Image = getImage(buildConf.Blockchain, "stable",buildConf.Blockchain)

		if !cpusEnabled {
			buildConf.Resources[0].Cpus = buildArr[offset]
			offset++
		}
		if !memoryEnabled {
			buildConf.Resources[0].Memory = buildArr[offset]
			offset++
		}

		if optionsFlag != nil {
			rawOptions, err := jsonRpcCall("get_params", []string{buildConf.Blockchain})
			if err != nil {
				util.PrintErrorFatal(err)
			}
			processedOptions, ok := rawOptions.([]interface{})
			if !ok {
				util.PrintStringError("Unexpected format for params")
				os.Exit(1)
			}
			buildConf.Params,err = processOptions(optionsFlag,processedOptions)
			if err != nil {
				util.PrintErrorFatal(err)
			}
		}else if len(paramsFile) != 0 {
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
			rawOptions, err := jsonRpcCall("get_params", []string{buildConf.Blockchain})
			if err != nil {
				util.PrintErrorFatal(err)
			}
			options, ok := rawOptions.([]interface{})
			if !ok {
				util.PrintStringError("Unexpected format for params")
				os.Exit(1)
			}

			scanner := bufio.NewScanner(os.Stdin)

			for i := 0; i < len(options); i++ {
				opt := options[i].([]interface{})
				if len(opt) != 2 {
					util.PrintStringError("Unexpected format for params")
					os.Exit(1)
				}

				key := opt[0].(string)
				key_type := opt[1].(string)

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
			if buildConf.Params == nil {
				buildConf.Params = map[string]interface{}{}
			}
			buildConf.Params["validators"] = validators
		}

		if buildConf.Blockchain == "eos" {
			if validators < 0 {
				buildConf.Nodes += 21
			} else {
				buildConf.Nodes += validators
			}
		}
		buildConf.Blockchain = strings.ToLower(buildConf.Blockchain)

		if filesFlag != nil{
			buildConf.Files = map[string]string{}
			for name,file := range filesFlag {
				data,err := ioutil.ReadFile(file)
				if err != nil{
					util.PrintErrorFatal(err)
				}
				buildConf.Files[name] = base64.StdEncoding.EncodeToString(data)

			}
		}
		
		if envFlag != nil {
			buildConf.Environments,err = processEnv(envFlag,buildConf.Nodes)
		}
		//fmt.Printf("%+v\n",buildConf)
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
		buildId,err := util.ReadStore(".in_progress_build_id")
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
	Short:  "Build a blockchain using previous configurations",
	Long: 	"\nBuild previous will recreate and deploy the previously built blockchain and specified number of nodes.\n",

	Run: func(cmd *cobra.Command, args []string) {

		prevBuild, err := getPreviousBuild()		
		if err != nil {
			util.PrintErrorFatal(err)
		}

		fmt.Println(prevBuild)
		if previousYesAll || util.YesNoPrompt("Build from previous?"){
			fmt.Println("building from previous configuration")
			build(prevBuild)
			removeSmartContracts()
			return
		}
	},
}

var buildStopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"halt", "cancel"},
	Short:   "Stops the current build",
	Long: "\nBuild stops the current building process.\n",

	Run: func(cmd *cobra.Command, args []string) {
		buildId,err := util.ReadStore(".in_progress_build_id")
		if err != nil || len(buildId) == 0 {
			fmt.Println("No inprogress build found. Use build command to deploy a blockchain.")
			os.Exit(1)
		}
		defer util.DeleteStore(".in_progress_build_id")
		jsonRpcCallAndPrint("stop_build", []string{string(buildId)})
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
	buildCmd.Flags().StringVarP(&imageFlag, "image", "i", "stable", "image tag")
	buildCmd.Flags().StringToStringVarP(&optionsFlag,"option","o",nil,"blockchain specific options")
	buildCmd.Flags().StringToStringVarP(&envFlag,"env","e",nil,"set environment variables for the nodes")
	buildCmd.Flags().StringToStringVarP(&filesFlag,"template","t",nil,"file templates")

	previousCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts. Evokes default parameters.")

	buildCmd.AddCommand(previousCmd, buildStopCmd, buildAttachCmd)
	RootCmd.AddCommand(buildCmd)
}
