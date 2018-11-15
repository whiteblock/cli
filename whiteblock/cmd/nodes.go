package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var listNodesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"nodes"},
	Short:   "List will show all nodes.",
	Long: `List will output all of the nodes in the current network.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		list := "docker ps -a | awk '{print $12}'"

		out, err := exec.Command("bash", "-c", list).Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
	},
}

func init() {
	RootCmd.AddCommand(listNodesCmd)
}
