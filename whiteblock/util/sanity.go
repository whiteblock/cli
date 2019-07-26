package util

import (
	"github.com/spf13/cobra"
)

func RequireFlags(cmd *cobra.Command, flags ...string) {
	for _, flag := range flags {
		if !cmd.Flags().Changed(flag) {
			FlagNotProvidedError(cmd, flag)
			return
		}
	}
}
