package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"github.com/whiteblock/cli/whiteblock/util"
)

var (
	mandir string
)

var manCmd = &cobra.Command{
	Hidden: true,
	Use:    "man",
	Short:  "Generate man pages for the Hugo CLI",
	Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		header := &doc.GenManHeader{
			Section: "1",
			Title:   "whiteblock",
		}

		cmd.Root().DisableAutoGenTag = true

		err := doc.GenManTree(cmd.Root(), header, mandir)
		if err != nil {
			util.PrintErrorFatal(err)
		}
	},
}

func init() {
	manCmd.Flags().StringVarP(&mandir, "dir", "d", "/tmp", "the directory to write the man pages.")
	RootCmd.AddCommand(manCmd)
}
