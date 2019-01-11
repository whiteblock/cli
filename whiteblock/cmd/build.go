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
	serverAddr 		string
	serversEnabled	bool
)

func writeFile(prevBuild string) {
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

func readFile() (string, error) {
	cwd := os.Getenv("HOME")
	b, err := ioutil.ReadFile(cwd + "/.config/whiteblock/previous_build.txt")
	if err != nil {
		fmt.Print(err)
	}
	return string(b), nil
}

func boolInput(input string) bool {
	output := false
	strings.ToLower(input)
	if input == "yes" || input == "on" || input == "true" || input == "y" {
		output = true
	}
	return output
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"init", "create"},
	Short:   "Build a blockchain using image and deploy nodes",
	Long: "Build will create and deploy a blockchain and the specified number of nodes."+
		  " Each node will be instantiated in its own container and will interact"+
		  " individually as a participant of the specified network.\n",

	Run: func(cmd *cobra.Command, args []string) {
		
		if len(args) != 0 {
			fmt.Println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		
		buildArr := make([]string, 0)
		paramArr := make([]string, 0)
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		bldcommand := "build"

		buildOpt := []string{"blockchain (" + blockchain + ")", "nodes (10)", "docker image (" + blockchain + ":latest)", 
								"cpus (empty for no limit)", "memory (empty for no limit)"}

		defOpt := []string{blockchain, "10", blockchain + ":latest", "", ""}

		if(serversEnabled){
			buildOpt = []string{"servers (" + server + ")", "blockchain (" + blockchain + ")", "nodes (10)", "docker image (" + blockchain + ":latest)", 
									"cpus (empty for no limit)", "memory (empty for no limit)"}
			defOpt = []string{fmt.Sprintf(server), blockchain, "10", blockchain + ":latest", "", ""}
		}
		

		

		scanner := bufio.NewScanner(os.Stdin)
		for i := 0; i < len(buildOpt); i++ {
			fmt.Print(buildOpt[i] + ": ")
			scanner.Scan()

			text := scanner.Text()
			if len(text) != 0 {
				buildArr = append(buildArr, text)
			} else {
				buildArr = append(buildArr, defOpt[i])
			}
		}

		if buildArr[0] == "[]" {
			fmt.Println("Invalid server. Please specify a server; none was given.")
			os.Exit(2)
		}
		
		var offset = 0
		if(serversEnabled){
			server = buildArr[offset]
			offset++
		}
		blockchain := buildArr[offset]
		offset++
		nodes := buildArr[offset]
		offset++
		image := buildArr[offset]
		offset++
		cpu := buildArr[offset]
		offset++
		memory := buildArr[offset]
		offset++

		fmt.Print("Use default parameters? (y/n) ")
		scanner.Scan()
		ask := scanner.Text()

		if ask != "n" {
		} else {
			getParamCommand := "get_params"
			bcparam := []byte(wsEmitListen(serverAddr, getParamCommand, blockchain))
			var paramlist []map[string]string

			json.Unmarshal(bcparam, &paramlist)

			scanner := bufio.NewScanner(os.Stdin)

			for i := 0; i < len(paramlist); i++ {
				for key, value := range paramlist[i] {
					fmt.Print(key, " ("+value+"): ")
					scanner.Scan()
					text := scanner.Text()
					if value == "string" {
						if len(text) != 0 {
							if fmt.Sprint(reflect.TypeOf(text)) != "string" {
								fmt.Println("bad type")
								os.Exit(2)
							}
							paramArr = append(paramArr, "\""+key+"\""+": "+"\""+text+"\"")
						} else {
							continue
						}
					} else if value == "[]string" {
						if len(text) != 0 {
							tmp := strings.Replace(text, " ", ",", -1)
							paramArr = append(paramArr, "\""+key+"\""+": "+"["+tmp+"]")
						} else {
							continue
						}
					} else if value == "int" {
						if len(text) != 0 {
							_, err := strconv.Atoi(text)
							if err != nil {
								fmt.Println("bad type")
								os.Exit(2)
							}
							paramArr = append(paramArr, "\""+key+"\""+": "+text)
						} else {
							continue
						}
					}
				}
			}
		}

		param := "{\"servers\":[" + server + "],\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + 
				 	"\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"},\"params\":{" + strings.Join(paramArr[:], ",") + "}}"
		stat := wsEmitListen(serverAddr, bldcommand, param)
		if stat == "" {
			writeFile(param)
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

		prevBuild, _ := readFile()

		if len(prevBuild) == 0 {
			fmt.Println("No previous build. Use build command to deploy a blockchain.")
		} else {
			fmt.Println(prevBuild)
			print("Build from previous? (y/n) ")
			scanner := bufio.NewScanner(os.Stdin)
			ask := scanner.Text()
			scanner.Scan()
			if ask != "n" {
				fmt.Println("building from previous configuration")
				wsEmitListen(serverAddr, bldcommand, prevBuild)
			} else {
				fmt.Println("Build cancelled.")
			}
		}

	},
}

func init() {
	buildCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	buildCmd.Flags().BoolVarP(&serversEnabled,"servers","s",false,"display server options")
	
	buildCmd.AddCommand(previousCmd)
	RootCmd.AddCommand(buildCmd)
}
