package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
	_, err := exec.Command("bash", "-c", "mkdir -p "+cwd+"/.config/whiteblock/").Output()
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

// func checkServ(server string) string {
// 	servList := make([]string, 0)
// 	servList = append(servList, "GUI to view stats and network information found here:")
// 	if strings.Contains(server, "1") {
// 		servList = append(servList, " 172.16.1.5:3000")
// 	}
// 	if strings.Contains(server, "2") {
// 		servList = append(servList, " 172.16.2.5:3000")
// 	}
// 	if strings.Contains(server, "3") {
// 		servList = append(servList, " 172.16.3.5:3000")
// 	}
// 	if strings.Contains(server, "4") {
// 		servList = append(servList, " 172.16.4.5:3000")
// 	}
// 	if strings.Contains(server, "5") {
// 		servList = append(servList, " 172.16.5.5:3000")
// 	}
// 	return strings.Join(servList, " ")
// }

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
		buildArr := make([]string, 0)
		paramArr := make([]string, 0)
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		bldcommand := "build"

		buildOpt := [6]string{"servers (default set to: [])", "blockchain (default set to: ethereum)", "nodes (default set to: 10)", "image (default set to: ethereum:latest)", "cpus (default set to: no limit)", "memory (default set to: no limit)"}
		defOpt := [6]string{"[]", "ethereum", "10", "ethereum:latest", "", ""}

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

		server := "[" + buildArr[0] + "]"
		blockchain := buildArr[1]
		nodes := buildArr[2]
		image := buildArr[3]
		cpu := buildArr[4]
		memory := buildArr[5]

		fmt.Print("Use default parameters? (y/n) ")
		scanner.Scan()
		ask := scanner.Text()

		if ask == "y" {
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
							i, err := strconv.Atoi(text)
							if err != nil {
								println("bad type")
								os.Exit(2)
							}
							paramArr = append(paramArr, "\""+key+"\""+": "+string(i))
						} else {
							continue
						}
					}
				}
			}
		}

		param := "{\"servers\":" + fmt.Sprintf("%s", server) + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"},\"params\":{" + strings.Join(paramArr[:], ",") + "}}"
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
			println("No previous build. Use build command to deploy a blockchain.")
		} else {
			wsEmitListen(serverAddr, bldcommand, prevBuild)
		}
	},
}

func init() {
	buildCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	buildCmd.AddCommand(previousCmd)
	RootCmd.AddCommand(buildCmd)
}
