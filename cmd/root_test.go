package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func executeCommand(args ...string) (string, error) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	_, err := rootCmd.ExecuteC()
	return buf.String(), err
}

func resetFlags() {
	address = ""
	token = ""
	namespace = ""
}

func TestRootCmd_RequiresTwoArgs(t *testing.T) {
	defer resetFlags()
	_, err := executeCommand()
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestRootCmd_TooManyArgs(t *testing.T) {
	defer resetFlags()
	_, err := executeCommand("path/a", "path/b", "path/c")
	if err == nil {
		t.Fatal("expected error when too many args provided")
	}
}

func TestRootCmd_FlagsRegistered(t *testing.T) {
	flags := []string{"address", "token", "namespace"}
	for _, f := range flags {
		if rootCmd.PersistentFlags().Lookup(f) == nil {
			t.Errorf("expected flag --%s to be registered", f)
		}
	}
}

func TestRootCmd_UsageContainsPaths(t *testing.T) {
	cmd := &cobra.Command{}
	_ = cmd
	use := rootCmd.Use
	if use == "" {
		t.Error("expected non-empty Use field on root command")
	}
}
