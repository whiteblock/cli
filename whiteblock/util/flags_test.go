package util

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestRequireFlags(t *testing.T) {
	// TODO test
}

func TestGetStringFlagValue(t *testing.T) {
	command := new(cobra.Command)
	command.Flags().String("test", "string", "blah")

	flag := "test"
	expected := "string"

	if !reflect.DeepEqual(GetStringFlagValue(command, flag), expected) {
		t.Error("return value of GetStringFlagValue")
	}
}

func TestGetIntFlagValue(t *testing.T) {
	command := new(cobra.Command)
	command.Flags().Int("test", 14, "blah")

	flag := "test"
	expected := 14

	if !reflect.DeepEqual(GetIntFlagValue(command, flag), expected) {
		t.Error("return value of GetStringFlagValue")
	}
}

func TestGetBoolFlagValue(t *testing.T) {
	command := new(cobra.Command)
	command.Flags().Bool("test", true, "blah")

	flag := "test"
	expected := true

	if !reflect.DeepEqual(GetBoolFlagValue(command, flag), expected) {
		t.Error("return value of GetStringFlagValue")
	}
}
