package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultdiff/internal/diff"
	"vaultdiff/internal/output"
	"vaultdiff/internal/vault"
)

// newRunCmd builds the cobra command that performs the actual diff between two Vault paths.
// It is registered as the primary action on the root command.
func newRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vaultdiff <left-path> <right-path>",
		Short: "Diff secrets between two Vault paths or namespaces",
		Long: `Compare secrets stored at two HashiCorp Vault KV paths.

Paths may include an optional namespace prefix separated by a pipe character:

  vaultdiff ns1|secret/app ns2|secret/app

Authentication is controlled via environment variables or flags (see --help).`,
		Args:    cobra.ExactArgs(2),
		RunE:    runDiff,
		SilenceUsage: true,
	}

	registerAuthFlags(cmd)
	registerFilterFlags(cmd)
	registerMaskFlags(cmd)

	cmd.Flags().StringP("format", "f", "text", `Output format: text, json, markdown`)
	cmd.Flags().BoolP("exit-code", "e", false, "Exit with code 1 if differences are found")

	return cmd
}

// runDiff is the core execution function invoked by the cobra command.
func runDiff(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	leftRaw := args[0]
	rightRaw := args[1]

	// Parse Vault path descriptors (namespace + mount + prefix).
	leftPath, err := vault.ParseVaultPath(leftRaw)
	if err != nil {
		return fmt.Errorf("invalid left path %q: %w", leftRaw, err)
	}
	rightPath, err := vault.ParseVaultPath(rightRaw)
	if err != nil {
		return fmt.Errorf("invalid right path %q: %w", rightRaw, err)
	}

	// Build Vault client and authenticate.
	authCfg, err := resolveAuthConfig(cmd)
	if err != nil {
		return fmt.Errorf("auth config: %w", err)
	}

	client, err := vault.NewClient(leftPath.Namespace)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	if err := vault.Authenticate(client, authCfg); err != nil {
		return fmt.Errorf("vault auth: %w", err)
	}

	// Resolve optional flags.
	filterOpts := resolveFilterOptions(cmd)
	maskOpts := resolveMaskOptions(cmd)

	format, _ := cmd.Flags().GetString("format")
	exitOnDiff, _ := cmd.Flags().GetBool("exit-code")

	// Read secrets from both sides.
	leftSecrets, err := vault.ReadAll(ctx, client, leftPath)
	if err != nil {
		return fmt.Errorf("reading left path: %w", err)
	}
	rightSecrets, err := vault.ReadAll(ctx, client, rightPath)
	if err != nil {
		return fmt.Errorf("reading right path: %w", err)
	}

	// Apply filters and masking.
	leftSecrets = vault.FilterSecrets(leftSecrets, filterOpts)
	rightSecrets = vault.FilterSecrets(rightSecrets, filterOpts)

	leftSecrets = vault.MaskSecrets(leftSecrets, maskOpts)
	rightSecrets = vault.MaskSecrets(rightSecrets, maskOpts)

	// Compute diff.
	results := diff.Compare(leftSecrets, rightSecrets)

	// Render output.
	renderer := output.NewRenderer(cmd.OutOrStdout())
	if err := renderer.Render(format, results); err != nil {
		return fmt.Errorf("render: %w", err)
	}

	// Optionally signal differences via exit code.
	if exitOnDiff && len(results) > 0 {
		os.Exit(1)
	}

	return nil
}
