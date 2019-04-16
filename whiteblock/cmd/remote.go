package cmd

import "github.com/spf13/cobra"

var remoteCmd = &cobra.Command{
	Aliases: []string{},
	Hidden:  true,
	Use:     "remote",
	Short:   "Remove all auth stored",
	Long:    "\nDeletes all stored auth\n",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("You are connected to %s\n", serverAddr)
	},
}

func init() {
	RootCmd.AddCommand(remoteCmd)
}
