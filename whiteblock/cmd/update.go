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

func BashExec(_cmd string) (string, error) {
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
	endpoint := fmt.Sprintf("https://storage.cloud.google.com/genesis-public/cli/%s/bin/%s/%s/whiteblock",
		branch, runtime.GOOS, runtime.GOARCH)
	binary, err := util.HttpRequest("GET", endpoint, "")
	if err != nil {
		util.PrintErrorFatal(err)
	}
	//let's find the where this binary is located.
	binaryLocation := os.Args[0]
	log.WithFields(log.Fields{"loc": binaryLocation}).Trace("got the binary location")
	if !strings.Contains(binaryLocation, "/") {
		binaryLocation, err = BashExec(fmt.Sprintf("which %s"))
		if err != nil {
			util.PrintErrorFatal(err)
		}
	} else {
		//convert to absolute path
		binaryLocation, err = filepath.Abs(binaryLocation)
		if err != nil {
			util.PrintErrorFatal(err)
		}
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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update the CLI",
	Long:  `Updates the cli binary`,

	Run: func(cmd *cobra.Command, args []string) {
		util.CheckArguments(cmd, args, 1, 1)
		var branch string
		switch args[0] {
		case "beta":
			fallthrough
		case "dev":
			branch = "dev"
		case "master":
			fallthrough
		case "stable":
			branch = "master"
		default:
			util.PrintErrorFatal("Invalid argument, specify either beta or stable")
		}
		switch runtime.GOOS {
		case "linux":
			handleUpdate(branch)
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
