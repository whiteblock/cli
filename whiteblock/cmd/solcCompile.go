package cmd

import (
	"fmt"

	solc "github.com/ethereum/go-ethereum/common/compiler"
	"github.com/spf13/cobra"
)

var (
	path     string
	filename string
)

func compile(path, filename string) {
	out, err := solc.CompileSolidity("solc", path+"/"+filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

var solcCMD = &cobra.Command{
	Use:   "contractcompile",
	Short: "Smart contract compiler.",
	Long: `
Compiles the specified smart contract.

	`,

	Run: func(cmd *cobra.Command, args []string) {
		compile(path, filename)
	},
}

func init() {
	solcCMD.LocalFlags().StringVarP(&path, "path", "p", "", "File path where the smart contract is located")
	solcCMD.LocalFlags().StringVarP(&filename, "filename", "f", "", "File name of the smart contract")

	RootCmd.AddCommand(solcCMD)
}
