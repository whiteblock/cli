// package cmd

// import (
// 	"log"

// 	"github.com/spf13/cobra"
// 	"github.com/spf13/cobra/doc"
// )

// func manp() {
// 	manCmd := &cobra.Command{
// 		Use:   "man",
// 		Short: "man page for information about the whiteblock cli application",
// 	}
// 	header := &doc.GenManHeader{
// 		Title:   "whiteblock",
// 		Section: "1",
// 	}
// 	err := doc.GenManTree(manCmd, header, "/tmp")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func init() {
// 	manp()
// }

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	mandir string
)

func newGenManCmd() {
	manCmd := &cobra.Command{
		Use:   "man",
		Short: "Generate man pages for the Hugo CLI",
		Long: `This command automatically generates up-to-date man pages of Hugo's
command-line interface.  By default, it creates the man page files
in the "man" directory under the current directory.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			header := &doc.GenManHeader{
				Section: "1",
				Title:   "whiteblock",
			}

			cmd.Root().DisableAutoGenTag = true

			err := doc.GenManTree(cmd.Root(), header, mandir)
			if err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}

	manCmd.Flags().StringVarP(&mandir, "dir", "d", "/tmp", "the directory to write the man pages.")

}
