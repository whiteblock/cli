package cmd

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/whiteblock/cli/whiteblock/util"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func bashExec(_cmd string) (string, error) {
	cmd := exec.Command("bash", "-c", _cmd)

	var resultsRaw bytes.Buffer

	cmd.Stdout = &resultsRaw
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return resultsRaw.String(), nil
}

func handleUpdate(branch string) {
	endpoint := fmt.Sprintf("https://storage.googleapis.com/genesis-public/cli/%s/bin/%s/%s/whiteblock",
		branch, runtime.GOOS, runtime.GOARCH)
	log.WithFields(log.Fields{"ep": endpoint}).Trace("fetching the binary data")
	binary, err := util.HttpRequest("GET", endpoint, "")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	log.WithFields(log.Fields{"size": len(binary)}).Trace("fetched the binary data")
	//let's find the where this binary is located.
	binaryLocation := os.Args[0]
	log.WithFields(log.Fields{"loc": binaryLocation}).Trace("got the binary location")
	if !strings.Contains(binaryLocation, "/") {
		binaryLocation, err = bashExec(fmt.Sprintf("which %s", binaryLocation))
	} else {
		binaryLocation, err = filepath.Abs(binaryLocation)
	}
	if err != nil {
		util.PrintErrorFatal(err)
	}

	binaryLocation, err = filepath.EvalSymlinks(binaryLocation)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	fi, err := os.Lstat(binaryLocation)
	if err != nil {
		util.PrintErrorFatal(err)
	}

	err = ioutil.WriteFile(binaryLocation+".tmp", binary, fi.Mode())
	if err != nil {
		util.PrintErrorFatal(err)
	}

	//swap,will be atomic on sane systems
	err = os.Rename(binaryLocation+".tmp", binaryLocation)
	if err != nil {
		util.PrintErrorFatal(err)
	}
	util.Print("whiteblock cli has been updated.")
}

func getUpdateBranch(args []string) string {
	if len(args) == 0 {
		options := []string{
			"master",
			"dev",
		}
		return options[util.OptionListPrompt("Which version of the cli would you like to use?", []string{
			"stable",
			"beta",
		})]
	}
	switch args[0] {
	case "beta":
		fallthrough
	case "dev":
		return "dev"
	case "master":
		fallthrough
	case "stable":
		return "master"
	default:
		util.PrintErrorFatal("Invalid argument, specify either beta or stable")
	}
	return ""
}

var updateCmd = &cobra.Command{
	Use:   "update [release]",
	Short: "Update the CLI",
	Long:  `Updates the cli binary to either the stable or beta release`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 0, 1)
		branch := getUpdateBranch(args)

		switch runtime.GOOS {
		case "linux":
			handleUpdate(branch)
		default:
			util.Print("sorry, your OS does not support easy updates")
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
