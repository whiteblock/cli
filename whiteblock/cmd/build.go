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
	serverAddr string
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
	Long: `
Build will create and deploy a blockchain and the specified number of nodes. Each node will be instantiated in its own container and will interact individually as a participant of the specified network.
	`,

	Run: func(cmd *cobra.Command, args []string) {

		if len(args) == 0 {
			buildArr := make([]string, 0)
			paramArr := make([]string, 0)
			serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
			bldcommand := "build"

			// buildOpt := [6]string{"servers (default set to: " + server + ")", "blockchain (default set to: " + blockchain + ")", "nodes (default set to: 10)", "image (default set to: " + blockchain + ":latest)", "cpus (default set to: no limit)", "memory (default set to: no limit)"}
			// defOpt := [6]string{fmt.Sprintf(server), blockchain, "10", blockchain + ":latest", "", ""}

			buildOpt := [5]string{"blockchain (default set to: " + blockchain + ")", "nodes (default set to: 10)", "image (default set to: " + blockchain + ":latest)", "cpus (default set to: no limit)", "memory (default set to: no limit)"}
			defOpt := [5]string{blockchain, "10", blockchain + ":latest", "", ""}

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
				println("Invalid server. Please specify a server; none was given.")
				os.Exit(2)
			}

			// server := "[" + buildArr[0] + "]"
			blockchain := buildArr[0]
			nodes := buildArr[1]
			image := buildArr[2]
			cpu := buildArr[3]
			memory := buildArr[4]

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
									println("bad type")
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
									println("bad type")
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

			param := "{\"servers\":" + server + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"},\"params\":{" + strings.Join(paramArr[:], ",") + "}}"
			stat := wsEmitListen(serverAddr, bldcommand, param)
			if stat == "" {
				writeFile(param)
			}
		} else {
			// param := "{\"servers\":" + args[0] + ",\"blockchain\":\"" + args[1] + "\",\"nodes\":" + args[2] + ",\"image\":\"" + args[3] + "\",\"resources\":{\"cpus\":\"" + args[4] + "\",\"memory\":\"" + args[5] + "\"},\"params\":{" + strings.Join(paramArr[:], ",") + "}}"
			// stat := wsEmitListen(serverAddr, bldcommand, param)
			// if stat == "" {
			// 	writeFile(param)
			// }
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
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
			println("No previous build. Use build command to deploy a blockchain.")
		} else {
			println(prevBuild)
			print("Build from previous? (y/n) ")
			scanner := bufio.NewScanner(os.Stdin)
			ask := scanner.Text()
			scanner.Scan()
			if ask != "n" {
				println("building from previous configuration")
				wsEmitListen(serverAddr, bldcommand, prevBuild)
			} else {
				println("Build cancelled.")
			}
		}

	},
}

func init() {
	buildCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	buildCmd.AddCommand(previousCmd)
	RootCmd.AddCommand(buildCmd)
}
