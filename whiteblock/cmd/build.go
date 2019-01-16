package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	serverAddr     string
	previousYesAll bool
	serversFlag    string
	blockchainFlag string
	nodesFlag      string
	cpusFlag       int
	memoryFlag     int
	paramsFile     string
	validators     int
)

type Config struct {
	Servers    []int
	Blockchain string
	Nodes      int
	Image      string
	Resources  struct {
		Cpus   string
		Memory string
	}
	Params interface{}
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
	/*
		fmt.Println(string(b))

		fmt.Println("get some of this:")

		fmt.Println(config.Servers)
		fmt.Println(config.Blockchain)
		fmt.Println(config.Nodes)
		fmt.Println(config.Image)
		fmt.Println(config.Resources.Cpus)
		fmt.Println(config.Resources.Memory)
	*/
	return b, nil
}

/*
func boolInput(input string) bool {
	output := false
	strings.ToLower(input)
	if input == "yes" || input == "on" || input == "true" || input == "y" {
		output = true
	}
	return output
}
*/

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
	//fmt.Println("server is: " + server)
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
		if len(nodesFlag) > 0 {
			nodesEnabled = true
		}
		if cpusFlag >= 0 {
			if cpusFlag == 0 {
				cpus = ""
			} else {
				cpus = strconv.Itoa(cpusFlag)
			}
			cpusEnabled = true
		}
		if memoryFlag >= 0 {
			if memoryFlag == 0 {
				memory = ""
			} else {
				memory = strconv.Itoa(memoryFlag)
			}
			memoryEnabled = true
		}

		if len(args) != 0 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}

		buildArr := make([]string, 0)
		paramArr := make(map[string]interface{})
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		bldcommand := "build"

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

		// fmt.Println(defaultBlockchain)
		// fmt.Println(server)
		// fmt.Println(defaultNodes)
		// fmt.Println(defaultImage)
		// fmt.Println(defaultCpus)
		// fmt.Println(defaultMemory)

		// fmt.Println(server)
		buildOpt := []string{}
		defOpt := []string{}
		allowEmpty := []bool{}

		// if !serversEnabled {
		// 	allowEmpty = []bool{false}
		// 	buildOpt = []string{
		// 		"servers" + tern((len(server) == 0), "", " ("+server+")"),
		// 	}
		// 	defOpt = append(defOpt, fmt.Sprintf(server))
		// }
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

		/*
			buildOpt = append(buildOpt, []string{
				"blockchain" + tern((len(defaultBlockchain) == 0), "", " ("+defaultBlockchain+")"),
				"nodes" + tern((defaultNodes == "0"), "", " ("+defaultNodes+")"),
				"docker image" + tern((len(defaultImage) == 0), "", " ("+defaultImage+")"),
				"cpus" + tern((len(defaultCpus) == 0), "(empty for no limit)", " ("+defaultCpus+")"),
				"memory" + tern((len(defaultMemory) == 0), "(empty for no limit)", " ("+defaultMemory+")"),
			}...)
		*/

		// defOpt = append(defOpt, []string{defaultCpus, defaultMemory}...)

		// allowEmpty = append(allowEmpty, []bool{true, true}...)

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
			nodes = buildArr[offset]
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
		params := "{}"

		if len(paramsFile) != 0 {
			f, err := os.Open(paramsFile)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			decoder := json.NewDecoder(f)
			decoder.UseNumber()
			err = decoder.Decode(&paramArr)
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
				getParamCommand := "get_params"
				bcparam := []byte(wsEmitListen(serverAddr, getParamCommand, blockchain))
				var paramlist []map[string]string

				json.Unmarshal(bcparam, &paramlist)

				scanner := bufio.NewScanner(os.Stdin)

				for i := 0; i < len(paramlist); i++ {
					for key, value := range paramlist[i] {
						fmt.Printf("%s (%s): ", key, value)
						scanner.Scan()
						text := scanner.Text()
						if len(text) == 0 {
							continue
						}
						switch value {
						case "string":
							if fmt.Sprint(reflect.TypeOf(text)) != "string" {
								fmt.Println("Entry must be a string")
								i--
								continue
							}
							paramArr[key] = "\"" + text + "\""
						case "[]string":
							paramArr[key] = "[" + strings.Replace(text, " ", ",", -1) + "]"
						case "int":
							_, err := strconv.ParseInt(text, 0, 64)
							if err != nil {
								fmt.Println("Entry must be an integer")
								i--
								continue
							}
							paramArr[key] = text
						}
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
			paramArr["validators"] = fmt.Sprintf("%d", validators)
		}

		params = "{"
		first := true
		for key, value := range paramArr {
			if first {
				first = false
			} else {
				params += ","
			}
			params += fmt.Sprintf("\"%s\""+":"+"%v", key, value)
		}
		params += "}"

		if blockchain == "eos" {
			ns,_ := strconv.Atoi(nodes)
			if validators < 0 {
				ns += 21
			}else{
				ns += validators
			}
			nodes = fmt.Sprintf("%d",ns)
		}

		param := "{\"servers\":[" + server + "],\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes +
			",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpus + "\",\"memory\":\"" + memory +
			"\"},\"params\":" + params + "}"
		stat := wsEmitListen(serverAddr, bldcommand, param)

		/*fmt.Println(blockchain)
		fmt.Println(server)
		fmt.Println(nodes)
		fmt.Println(image)
		fmt.Println(cpus)
		fmt.Println(memory)

		fmt.Println(bldcommand)
		fmt.Println(param)*/

		if stat == "" {
			writePrevCmdFile(param)
			writeConfigFile(param)
		}

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
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		bldcommand := "build"
		prevBuild, _ := readPrevCmdFile()

		if len(prevBuild) == 0 {
			fmt.Println("No previous build. Use build command to deploy a blockchain.")
			return
		}

		fmt.Println(prevBuild)
		if previousYesAll {
			fmt.Println("building from previous configuration")
			wsEmitListen(serverAddr, bldcommand, prevBuild)
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
				wsEmitListen(serverAddr, bldcommand, prevBuild)
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
	buildCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts")
	buildCmd.Flags().StringVarP(&blockchainFlag, "blockchain", "b", "", "specify blockchain")
	buildCmd.Flags().StringVarP(&nodesFlag, "nodes", "n", "", "specify number of nodes")
	buildCmd.Flags().IntVarP(&cpusFlag, "cpus", "c", 0, "specify number of cpus")
	buildCmd.Flags().IntVarP(&memoryFlag, "memory", "m", 0, "specify memory allocated")
	buildCmd.Flags().StringVarP(&paramsFile, "file", "f", "", "parameters file")
	buildCmd.Flags().IntVarP(&validators, "validators", "v", -1, "set the number of validators")

	previousCmd.Flags().BoolVarP(&previousYesAll, "yes", "y", false, "Yes to all prompts")

	buildCmd.AddCommand(previousCmd)
	RootCmd.AddCommand(buildCmd)
}
