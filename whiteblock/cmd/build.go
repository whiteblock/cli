package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	serverAddr     string
	previousYesAll bool
	serversFlag    string
	blockchainFlag string
	nodesFlag      int
	cpusFlag       string
	memoryFlag     string
	paramsFile     string
	validators     int
)

type Config struct {
	Servers    	[]int					`json:"servers"`
	Blockchain 	string					`json:"blockchain"`
	Nodes      	int						`json:"nodes"`
	Image      	string					`json:"image"`
	Resources 	Resources				`json:"resources"`
	Params 		map[string]interface{}	`json:"params"`
}

type Resources struct {
	Cpus   string	`json:"cpus"`
	Memory string	`json:"memory"`
}

func writePrevCmdFile(prevBuild string) {
	cwd := os.Getenv("HOME")
	err := os.MkdirAll(cwd+"/.config/whiteblock/", 0755)
	if err != nil {
		log.Fatalf("could not create directory: %s", err)
	}

	file, err := os.Create(cwd + "/.config/whiteblock/previous_build.txt")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	_, err = file.WriteString(prevBuild)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func readPrevCmdFile() (string, error) {
	cwd := os.Getenv("HOME")
	b, err := ioutil.ReadFile(cwd + "/.config/whiteblock/previous_build.txt")
	if err != nil {
		//fmt.Print(err)
		return "",err
	}
	return string(b), nil
}

func writeConfigFile(configFile string) {
	cwd := os.Getenv("HOME")
	err := os.MkdirAll(cwd+"/.config/whiteblock/", 0755)
	if err != nil {
		log.Fatalf("could not create directory: %s", err)
	}

	file, err := os.Create(cwd + "/.config/whiteblock/config.json")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer file.Close() // Make sure to close the file when you're done

	_, err = file.WriteString(configFile)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
}

func readConfigFile() ([]byte, error) {
	cwd := os.Getenv("HOME")
	b, err := ioutil.ReadFile(cwd + "/.config/whiteblock/config.json")
	if err != nil {
		//fmt.Print(err)
	}
	return b, nil
}

func build(buildConfig Config) {
	buildReply,err := jsonRpcCall("build",buildConfig)
	if err != nil{
		PrintErrorFatal(err)
	}
	fmt.Printf("%v\n",buildReply)
	buildListener()
}


func getServer() string {
	idList := make([]string, 0)
	getServerAddr := serverAddr
	command := "get_servers"
	serverResults := []byte(wsEmitListen(getServerAddr, command, ""))
	var result map[string]Server
	err := json.Unmarshal(serverResults, &result)
	if err != nil {
		panic(err)
	}

	serverID := 0
	for _, v := range result {
		serverID = v.Id
		//move this and take out break statement if instance has multiple servers
		idList = append(idList, fmt.Sprintf("%d", serverID))
		break
	}

	server = strings.Join(idList, ",")
	return server
}

func tern(exp bool, res1 string, res2 string) string {
	if exp {
		return res1
	}
	return res2
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: "Build will create and deploy a blockchain and the specified number of nodes." +
		" Each node will be instantiated in its own container and will interact" +
		" individually as a participant of the specified network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		serversEnabled := false
		blockchainEnabled := false
		nodesEnabled := false
		cpusEnabled := false
		memoryEnabled := false
		if len(serversFlag) > 0 {
			serversEnabled = true
		}
		if len(blockchainFlag) > 0 {
			blockchainEnabled = true
		}
		if nodesFlag > 0 {
			nodesEnabled = true
		}
		if len(cpusFlag) != 0 {
			if cpusFlag == "0" {
				cpusFlag = ""
			} else {
				cpus = cpusFlag
			}
			cpusEnabled = true
		}
		if len(memoryFlag) > 0 {
			if memoryFlag == "0" {
				memoryFlag = ""
			} else {
				memory = memoryFlag
			}
			memoryEnabled = true
		}

		if len(args) != 0 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		buildArr := make([]string, 0)
		params := make(map[string]interface{})
		

		configFile, err := readConfigFile()
		if err != nil {
			panic(err)
		}

		var config Config
		json.Unmarshal(configFile, &config)
		defaultBlockchain := string(config.Blockchain)
		defaultNodes := strconv.Itoa(config.Nodes)
		//defaultImage := string(config.Image)
		defaultCpus := string(config.Resources.Cpus)
		defaultMemory := string(config.Resources.Memory)

		buildOpt := []string{}
		defOpt := []string{}
		allowEmpty := []bool{}

		/*
			if !serversEnabled {
				allowEmpty = []bool{false}
				buildOpt = []string{
					"servers" + tern((len(server) == 0), "", " ("+server+")"),
				}
				defOpt = append(defOpt, fmt.Sprintf(server))
			}
		*/
		if !blockchainEnabled {
			allowEmpty = append(allowEmpty, false)
			buildOpt = append(buildOpt, "blockchain"+tern((len(defaultBlockchain) == 0), "", " ("+defaultBlockchain+")"))
			defOpt = append(defOpt, fmt.Sprintf(defaultBlockchain))
		}
		if !nodesEnabled {
			allowEmpty = append(allowEmpty, false)
			buildOpt = append(buildOpt, "nodes"+tern((defaultNodes == "0"), "", " ("+defaultNodes+")"))
			defOpt = append(defOpt, fmt.Sprintf(defaultNodes))
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

		if len(serversFlag) == 0 {
			server = string(getServer())
		}

		var offset = 0
		if serversEnabled {
			server = serversFlag
		}
		if blockchainEnabled {
			blockchain = blockchainFlag
		} else {
			blockchain = buildArr[offset]
			offset++
		}

		if nodesEnabled {
			nodes = nodesFlag
		} else {
			nodes,err = strconv.Atoi(buildArr[offset])
			if err != nil{
				InvalidInteger("nodes",buildArr[offset],true)
			}
			offset++
		}

		image := "gcr.io/whiteblock/" + blockchain + ":master"
		// image := blockchain
		if !cpusEnabled {
			cpus = buildArr[offset]
			offset++
		}
		if !memoryEnabled {
			memory = buildArr[offset]
			offset++
		}

		if len(paramsFile) != 0 {
			f, err := os.Open(paramsFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			decoder := json.NewDecoder(f)
			decoder.UseNumber()
			err = decoder.Decode(&params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else if !previousYesAll {
		Params:
			fmt.Print("Use default parameters? (y/n) ")
			scanner.Scan()
			ask := scanner.Text()
			ask = strings.Trim(ask, "\n\t\r\v ")

			switch ask {
			case "n":
				fallthrough
			case "no":
				rawOptions,err := jsonRpcCall("get_params",nil)
				if err != nil{
					fmt.Println(err)
					os.Exit(1)
				}
				options,ok := rawOptions.([][]string)
				if !ok {
					fmt.Println("Unexpected format for params")
					os.Exit(1)
				}

				scanner := bufio.NewScanner(os.Stdin)

				for i := 0;i < len(options);i++ {
					opt := options[0]
					if len(opt) != 2 {
						fmt.Println("Unexpected format for params")
						os.Exit(1)
					}

					key := opt[0]
					key_type := opt[1]

					fmt.Printf("%s (%s): ", key, key_type)
					scanner.Scan()
					text := scanner.Text()
					if len(text) == 0 {
						continue
					}
					switch key_type {
						case "string":
							//needs to have filtering
							params[key] = text
						case "[]string":
							preprocessed := strings.Replace(text, " ", ",", -1)
							params[key] = strings.Split(preprocessed,",")
						case "int":
							val, err := strconv.ParseInt(text, 0, 64)
							if err != nil {
								InvalidInteger(key,text,false)
								i--
								continue
							}
							params[key] = val
					}
				}
				
				
			case "y":
				fallthrough
			case "yes":
			default:
				fmt.Println("Unknown Option")
				goto Params
			}
		}
		if validators >= 0 {
			params["validators"] =  validators
		}

		if blockchain == "eos" {
			if validators < 0 {
				nodes += 21
			} else {
				nodes += validators
			}
		}

		serverNum,_ := strconv.Atoi(server)

		buildConfig := Config{
			Servers: []int{serverNum},
			Blockchain: blockchain,
			Nodes: nodes,
			Image: image,
			Resources: Resources{
				Cpus:cpus,
				Memory:memory,
			},
			Params:params,
		}

		build(buildConfig)

		param,err := json.Marshal(buildConfig)
		writePrevCmdFile(string(param))
		writeConfigFile(string(param))
	},
}

var previousCmd = &cobra.Command{
	Use:     "previous",
	Aliases: []string{"prev"},
	Short:   "Build a blockchain using previous configurations",
	Long: `
Build previous will recreate and deploy the previously built blockchain and specified number of nodes.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		rawPrevBuild, _ := readPrevCmdFile()
		if len(rawPrevBuild) == 0 {
			fmt.Println("No previous build. Use build command to deploy a blockchain.")
			os.Exit(1)
		}
		var prevBuild Config
		err :=  json.Unmarshal([]byte(rawPrevBuild),&prevBuild)
		if err != nil{
			log.Println(err)
			os.Exit(1)
			//PrintErrorFatal(err)
		}

		fmt.Println(prevBuild)
		if previousYesAll {
			fmt.Println("building from previous configuration")
			build(prevBuild)
			return
		}

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("Build from previous? (y/n) ")
			scanner.Scan()
			ask := scanner.Text()
			ask = strings.Trim(ask, "\n\t\r\v ")

			switch ask {
				case "y":
					fallthrough
				case "yes":
					fmt.Println("building from previous configuration")
					build(prevBuild)
					return
				case "n":
					fallthrough
				case "no":
					fmt.Println("Build cancelled.")
					return
				default:
					fmt.Println("Unknown Option " + ask)
			}
		}

	},
}

func init() {
	buildCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	buildCmd.Flags().StringVarP(&serversFlag, "servers", "s", "", "display server options")
	buildCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts. Evokes default parameters.")
	buildCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	buildCmd.Flags().IntVarP(&nodesFlag, "nodes", "n", 0, "specify number of nodes")
	buildCmd.Flags().StringVarP(&cpusFlag, "cpus", "c", "", "specify number of cpus")
	buildCmd.Flags().StringVarP(&memoryFlag, "memory", "m", "", "specify memory allocated")
	buildCmd.Flags().StringVarP(&paramsFile, "file", "f", "", "parameters file")
	buildCmd.Flags().IntVarP(&validators, "validators", "v", -1, "set the number of validators")

	previousCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	previousCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts. Evokes default parameters.")

	buildCmd.AddCommand(previousCmd)
	RootCmd.AddCommand(buildCmd)
}
