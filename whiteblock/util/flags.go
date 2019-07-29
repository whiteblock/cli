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

func GetStringFlagValue(cmd *cobra.Command, flag string) string {
	out, err := cmd.Flags().GetString(flag)
	if err != nil {
		PrintErrorFatal(err)
	}
	return out
}

func GetIntFlagValue(cmd *cobra.Command, flag string) int {
	out, err := cmd.Flags().GetInt(flag)
	if err != nil {
		PrintErrorFatal(err)
	}
	return out
}

func GetBoolFlagValue(cmd *cobra.Command, flag string) bool {
	out, err := cmd.Flags().GetBool(flag)
	if err != nil {
		PrintErrorFatal(err)
	}
	return out
}
