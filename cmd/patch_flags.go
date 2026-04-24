package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

// registerPatchFlags adds patch-related flags to the given command.
func registerPatchFlags(cmd *cobra.Command) {
	cmd.Flags().StringArray("patch-set", nil, "Set a key: path/key=value")
	cmd.Flags().StringArray("patch-delete", nil, "Delete a key: path/key")
	cmd.Flags().StringArray("patch-rename", nil, "Rename a key: path/oldkey=newkey")
	cmd.Flags().Bool("patch-dry-run", false, "Preview patch operations without applying them")
}

// resolvePatchOptions builds PatchOptions from command-line flags.
func resolvePatchOptions(cmd *cobra.Command) (vault.PatchOptions, error) {
	var ops []vault.PatchOperation

	sets, _ := cmd.Flags().GetStringArray("patch-set")
	for _, s := range sets {
		path, key, value, err := parsePathKeyValue(s)
		if err != nil {
			return vault.PatchOptions{}, fmt.Errorf("--patch-set %q: %w", s, err)
		}
		ops = append(ops, vault.PatchOperation{Op: "set", Path: path, Key: key, Value: value})
	}

	deletes, _ := cmd.Flags().GetStringArray("patch-delete")
	for _, d := range deletes {
		path, key, err := parsePathKey(d)
		if err != nil {
			return vault.PatchOptions{}, fmt.Errorf("--patch-delete %q: %w", d, err)
		}
		ops = append(ops, vault.PatchOperation{Op: "delete", Path: path, Key: key})
	}

	renames, _ := cmd.Flags().GetStringArray("patch-rename")
	for _, r := range renames {
		path, oldKey, newKey, err := parsePathKeyValue(r)
		if err != nil {
			return vault.PatchOptions{}, fmt.Errorf("--patch-rename %q: %w", r, err)
		}
		ops = append(ops, vault.PatchOperation{Op: "rename", Path: path, Key: oldKey, NewKey: newKey})
	}

	dryRun, _ := cmd.Flags().GetBool("patch-dry-run")
	return vault.PatchOptions{Operations: ops, DryRun: dryRun}, nil
}

// parsePathKeyValue splits "some/path/key=value" into ("some/path", "key", "value", nil).
func parsePathKeyValue(s string) (string, string, string, error) {
	eq := strings.LastIndex(s, "=")
	if eq < 0 {
		return "", "", "", fmt.Errorf("expected path/key=value")
	}
	left, value := s[:eq], s[eq+1:]
	path, key, err := parsePathKey(left)
	return path, key, value, err
}

// parsePathKey splits "some/path/key" into ("some/path", "key", nil).
func parsePathKey(s string) (string, string, error) {
	idx := strings.LastIndex(s, "/")
	if idx < 0 {
		return "", "", fmt.Errorf("expected path/key")
	}
	return s[:idx], s[idx+1:], nil
}
