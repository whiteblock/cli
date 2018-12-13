package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type List struct {
	Results []struct {
		Series []struct {
			Columns []string        `json:"columns"`
			Values  [][]interface{} `json:"values"`
		}
	} `json:"results"`
}

var sysCMD = &cobra.Command{
	Use:   "sys <command>",
	Short: "Run SYS commands.",
	Long: `
Sys will allow the user to get infromation and run SYS commands.
`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock sys -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var sysTestCMD = &cobra.Command{
	Use:   "test <command>",
	Short: "SYS test commands.",
	Long: `
Sys test will allow the user to get infromation and run SYS tests.
`,

	Run: func(cmd *cobra.Command, args []string) {
		out, err := exec.Command("bash", "-c", "./whiteblock sys test -h").Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
		println("\nNo command given. Please choose a command from the list above.\n")
		os.Exit(1)
	},
}

var testStartCMD = &cobra.Command{
	Use:   "start <wait time> <min complete percent> <number of tx>",
	Short: "Starts propagation test.",
	Long: `
Sys test start will start the propagation test. It will wait for the signal start time, have nodes send messages at the same time, and require to wait a minimum amount of time then check receivers with a completion rate of minimum completion percentage. 

Format: <wait time> <min complete percent> <number of tx>
Params: Time in seconds, percentage, number of transactions

`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "sys::start_test"
		if len(args) != 3 {
			out, err := exec.Command("bash", "-c", "./whiteblock sys test start -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}
		param := "{\"waitTime\":" + args[0] + ",\"minCompletePercent\":" + args[1] + ",\"numberOfTransactions\":" + args[2] + "}"
		// println(command)
		// println(param)
		wsEmitListen(serverAddr, command, param)
	},
}

var testResultsCMD = &cobra.Command{
	Use:   "results <test number>",
	Short: "Get results from a previous test.",
	Long: `
Sys test results pulls data from a previous test or tests and outputs as csv.

Format: <test number>
Params: Test number

	`,

	Run: func(cmd *cobra.Command, args []string) {
		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		command := "sys::get_recent_test_results"
		if len(args) != 1 {
			out, err := exec.Command("bash", "-c", "./whiteblock sys test results -h").Output()
			if err != nil {
				os.Exit(1)
			}
			fmt.Printf("%s", out)
			println("\nError: Invalid number of arguments given\n")
			os.Exit(1)
		}
		results := []byte(wsEmitListen(serverAddr, command, args[0]))
		var result map[string]interface{}
		err := json.Unmarshal(results, &result)
		if err != nil {
			panic(err)
		}

		var l List
		json.Unmarshal(results, &l)
		r := l.Results
		rc := r[0]
		s := rc.Series
		sc := s[0]
		c := sc.Columns
		v := sc.Values[0]

		for i := 0; i < len(c); i++ {
			fmt.Println("\t" + c[i] + ": " + fmt.Sprint(v[i]) + " type is: " + fmt.Sprint(reflect.TypeOf(v[i])))
		}
	},
}

var testPinnacleCMD = &cobra.Command{
	Use:   "pinnacle <wait time> <min complete percent>",
	Short: "Run the pinnacle test series.",
	Long: `
Sys test pinnacle will run the propagation test.

Format: <test number>
Params: Test number

	`,

	Run: func(cmd *cobra.Command, args []string) {
		var pinArgs []string

		// println("connecting to: " + serverAddr)

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		xS := "1000"
		y := 500
		z := 0

		buildCmd.Run = func(cmd *cobra.Command, args []string) {
			buildArr := make([]string, 0)
			paramArr := make([]string, 0)
			serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
			bldcommand := "build"

			buildOpt := [6]string{"servers (required server number)", "blockchain (default set to: ethereum)", "nodes (default set to: 10)", "image (default set to: ethereum:latest)", "cpus (default set to: no limit)", "memory (default set to: no limit)"}
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

			if buildArr[0] == "[]" {
				println("Invalid server. Please specify a server; none was given.")
				os.Exit(2)
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

			param := "{\"servers\":" + fmt.Sprintf("%s", server) + ",\"blockchain\":\"" + blockchain + "\",\"nodes\":" + nodes + ",\"image\":\"" + image + "\",\"resources\":{\"cpus\":\"" + cpu + "\",\"memory\":\"" + memory + "\"},\"params\":{" + strings.Join(paramArr[:], ",") + "}}"
			stat := wsEmitListen(serverAddr, bldcommand, param)
			if stat == "" {
				writeFile(param)
			}
		}

		pinArgs = append(pinArgs, args[0])
		pinArgs = append(pinArgs, args[1])
		pinArgs = append(pinArgs, xS)
		testStartCMD.Run(testStartCMD, pinArgs[:])

		command := "sys::get_recent_test_results"
		results := []byte(wsEmitListen(serverAddr, command, args[0]))
		var result map[string]interface{}
		err := json.Unmarshal(results, &result)
		if err != nil {
			panic(err)
		}

		var l List
		json.Unmarshal(results, &l)
		r := l.Results
		rc := r[0]
		s := rc.Series
		sc := s[0]
		c := sc.Columns

		avgTestTime := c[1]

		for i := 0; i < len(pinArgs); i++ {
			println(pinArgs[i])
		}

		xI, err := strconv.Atoi(xS)
		if err != nil {
			panic(err)
		}

		for {
			command := "sys::get_recent_test_results"
			results := []byte(wsEmitListen(serverAddr, command, args[0]))
			var result map[string]interface{}
			err := json.Unmarshal(results, &result)
			if err != nil {
				panic(err)
			}

			var l List
			json.Unmarshal(results, &l)
			r := l.Results
			rc := r[0]
			s := rc.Series
			sc := s[0]
			c := sc.Columns

			if c[1] != avgTestTime {
				avgTestTime = c[1]
				if c[9] == "0" {
					if y > 50 {
						xI = (xI - y)
						xS = string(xI)
						z = xI
					} else {
						break
					}
				} else if c[9] == "1" {
					xI = (xI + y)
					y = y / 2
					xS = string(xI)
					z = xI
				}
			}
			pinArgs[2] = xS
			previousCmd.Execute()
			testStartCMD.Run(testStartCMD, pinArgs[:])
		}
		println(z)
	},
}

func init() {
	testStartCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	testResultsCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	sysTestCMD.AddCommand(testStartCMD, testResultsCMD, testPinnacleCMD)
	sysCMD.AddCommand(sysTestCMD)
	RootCmd.AddCommand(sysCMD)
}
