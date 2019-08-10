package build

import (
	"fmt"
	"github.com/whiteblock/cli/whiteblock/util"
)

type Status struct {
	Error    map[string]interface{} `json:"error"`
	Progress float64                `json:"progress"`
	Stage    string                 `json:"stage"`
	Frozen   bool                   `json:"frozen"`
}

func (status Status) Print() bool {
	if status.Frozen {
		fmt.Printf("\nBuild is currently frozen. Press Ctrl-\\ to drop into console. Run 'whiteblock build unfreeze' to resume. \r")
	} else if status.Error != nil {
		fmt.Println() //move to the next line
		util.PrintErrorFatal(status.Error["what"])
	} else if status.Progress == 0.0 {
		fmt.Printf("Sending build context to Whiteblock\r")
	} else if status.Progress == 100.0 {
		fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\r", "Build", status.Progress)
		util.Print("\a")
		return true
	} else {
		fmt.Printf("\033[1m\033[K\033[31m%s\033[0m\t%f%% completed\r", status.Stage, status.Progress)
	}
	return false
}
