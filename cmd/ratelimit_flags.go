package cmd

import (
	"github.com/spf13/cobra"

	"github.com/your-org/vaultdiff/internal/vault"
)

// registerRateLimitFlags attaches rate-limit flags to the given command.
func registerRateLimitFlags(cmd *cobra.Command) {
	cmd.Flags().Float64(
		"rate-limit",
		10,
		"Maximum sustained Vault API requests per second (0 = unlimited)",
	)
	cmd.Flags().Float64(
		"rate-burst",
		20,
		"Maximum burst size for Vault API requests",
	)
}

// resolveRateLimitOptions builds RateLimitOptions from the command flags.
// When rate-limit is 0, a high ceiling is used to approximate "unlimited".
func resolveRateLimitOptions(cmd *cobra.Command) vault.RateLimitOptions {
	rps, _ := cmd.Flags().GetFloat64("rate-limit")
	burst, _ := cmd.Flags().GetFloat64("rate-burst")

	if rps <= 0 {
		// Effectively unlimited — use a very large value.
		rps = 1_000_000
	}
	if burst <= 0 {
		burst = rps
	}
	return vault.RateLimitOptions{
		RequestsPerSecond: rps,
		Burst:             burst,
	}
}
