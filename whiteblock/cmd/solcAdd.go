package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	addPath     string
	addFilename string
)

var addscCMD = &cobra.Command{
	Use:   "contractadd",
	Short: "Add a smart contract.",
	Long: `Adds the specified smart contract into the /Downloads folder.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		cp := "cp " + fmt.Sprintf(addPath) + "/" + fmt.Sprintf(addFilename)

		out, err := exec.Command("bash", "-c", cp).Output()
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", out)
	},
}

func init() {
	addscCMD.Flags().StringVarP(&addPath, "path", "p", "", "File path where the smart contract is located")
	addscCMD.Flags().StringVarP(&addFilename, "filename", "f", "", "File name of the smart contract")

	RootCmd.AddCommand(addscCMD)
}
