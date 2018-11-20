package cmd

import (
	"github.com/spf13/cobra"
)

var listNodesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"nodes"},
	Short:   "List will show all nodes.",
	Long: `List will output all of the nodes in the current network.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		// 	list := "docker ps -a | awk '{print $12}'"

		// 	fmt.Println(list)

		// 	out, err := exec.Command("bash", "-c", list).Output()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	fmt.Printf("%s", out)
		// },

		serverAddr = "ws://" + serverAddr + "/socket.io/?EIO=3&transport=websocket"

		wsGetNodes(serverAddr)
	},
}

func init() {
	listNodesCmd.Flags().StringVarP(&serverAddr, "server-addr", "a", "localhost:5000", "server address with port 5000")

	RootCmd.AddCommand(listNodesCmd)
}
