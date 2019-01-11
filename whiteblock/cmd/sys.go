package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

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
		println("\nNo command given. Please choose a command from the list above.\n")
		cmd.Help()
		return
	},
}

var sysTestCMD = &cobra.Command{
	Use:   "test <command>",
	Short: "SYS test commands.",
	Long: `
Sys test will allow the user to get infromation and run SYS tests.
`,

	Run: func(cmd *cobra.Command, args []string) {
		println("\nNo command given. Please choose a command from the list above.\n")
		cmd.Help()
		return
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
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
		}
		param := "{\"waitTime\":" + args[0] + ",\"minCompletePercent\":" + args[1] + ",\"numberOfTransactions\":" + args[2] + "}"
		// println(command)
		// println(param)
		wsEmitListen(serverAddr, command, param)
		return
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
			println("\nError: Invalid number of arguments given\n")
			cmd.Help()
			return
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
		var wg sync.WaitGroup
		var emptyArgs []string
		var pinArgs []string
		var avgTestTime string

		xS := "200"
		y := 500
		z := 0

		buildCmd.Run(buildCmd, emptyArgs)

		time.Sleep(5 * time.Second)

		wg.Add(1)

		pinArgs = append(pinArgs, args[0])
		pinArgs = append(pinArgs, args[1])
		pinArgs = append(pinArgs, xS)

		for i := 0; i < len(pinArgs); i++ {
			println(pinArgs[i])
		}

		// sysStartServerAddr := "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"
		sysStartServerAddr := serverAddr
		pinnacleServerAddr := serverAddr

		go func() {
			param := "{\"waitTime\":" + pinArgs[0] + ",\"minCompletePercent\":" + pinArgs[1] + ",\"numberOfTransactions\":" + pinArgs[2] + "}"
			println(sysStartServerAddr)
			wsEmitListen(sysStartServerAddr, "sys::start_test", param)
			wg.Done()
		}()

		println("test start completed")
		println("waiting to get results")

		for {
			time.Sleep(10 * time.Second)
			runningTest := wsEmitListen(sysStartServerAddr, "state::is_running", "")
			println(runningTest)
			println(wsEmitListen(sysStartServerAddr, "state::what_is_running", ""))
			// how long does it take to sync sys::get_recent_test_results. there is an index error because there is no data.
			if runningTest == "false" {
				// println("getting test results")
				command := "sys::get_recent_test_results"
				results := []byte(wsEmitListen(sysStartServerAddr, command, args[0]))
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

				avgTestTime = fmt.Sprint(v[1])

				for j := 0; j < len(v); j++ {
					println(string(j) + " -- " + fmt.Sprint(c[j]) + " : " + fmt.Sprint(v[j]))
				}
				println("v[9] = " + fmt.Sprint(v[9]))
				break
			}
		}

		println(avgTestTime)

		xI, err := strconv.Atoi(xS)
		if err != nil {
			panic(err)
		}

		for {
			println("going inside pinnacle loops")

			println(pinnacleServerAddr)

			println("chaging pinArgs[2] to -> " + xS)
			pinArgs[2] = xS
			for k := 0; k < len(pinArgs); k++ {
				println(pinArgs[k])
			}

			println("build prev command")
			prevBuild, _ := readPrevCmdFile()
			println(prevBuild)
			wsEmitListen(pinnacleServerAddr, "build", prevBuild)
			time.Sleep(5 * time.Second)

			go func() {
				println("running test again")
				param := "{\"waitTime\":" + pinArgs[0] + ",\"minCompletePercent\":" + pinArgs[1] + ",\"numberOfTransactions\":" + pinArgs[2] + "}"
				println(pinnacleServerAddr)
				wsEmitListen(pinnacleServerAddr, "sys::start_test", param)
			}()

			for {
				time.Sleep(10 * time.Second)
				runningTest := wsEmitListen(sysStartServerAddr, "state::is_running", "")
				println(runningTest)
				println(wsEmitListen(sysStartServerAddr, "state::what_is_running", ""))
				// how long does it take to sync sys::get_recent_test_results. there is an index error because there is no data.
				if runningTest == "false" {
					println("getting test results")
					command := "sys::get_recent_test_results"
					results := []byte(wsEmitListen(sysStartServerAddr, command, args[0]))
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

					for j := 0; j < len(v); j++ {
						println(string(j) + " -- " + fmt.Sprint(c[j]) + " : " + fmt.Sprint(v[j]))
					}
					println("v[9] = " + fmt.Sprint(v[9]))

					//unit test this

					println(fmt.Sprint(v[1]) + " =?= " + fmt.Sprint(avgTestTime))
					if fmt.Sprint(v[1]) != avgTestTime {
						avgTestTime = fmt.Sprint(v[1])
						if fmt.Sprint(v[9]) == "0" {
							println("0 success code")
							if y > 50 {
								xI = (xI - y)
								xS = fmt.Sprintf("%d", xI)
								z = xI

								fmt.Println(xI)
								fmt.Println(y)
								fmt.Println(xS)
								println("getting out of 0 y>50")
								break
							} else {
								fmt.Println(xI)
								fmt.Println(y)
								fmt.Println(xS)
								println("getting out of 0 y<50")
								break
							}
						} else if fmt.Sprint(v[9]) == "1" {
							println("1 success code")
							xI = (xI + y)
							y = y / 2
							xS = fmt.Sprintf("%d", xI)
							z = xI

							fmt.Println(xI)
							fmt.Println(y)
							fmt.Println(xS)
							println("getting out of 1")
							break
						}
					} else {
						println("avgTestTime is equal")
						break
					}
				}
			}
			if y <= 50 {
				break
			}
			// for j := 0; j < len(v); j++ {
			// 	println(string(j) + " -- " + " : " + fmt.Sprint(v[j]))
			// }

		}
		println(z)
	},
}

func init() {
	testStartCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")
	// testResultsCMD.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	sysTestCMD.AddCommand(testStartCMD, testResultsCMD, testPinnacleCMD)
	sysCMD.AddCommand(sysTestCMD)
	RootCmd.AddCommand(sysCMD)
}
