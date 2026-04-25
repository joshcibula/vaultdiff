package cmd

import (
	"github.com/spf13/cobra"
	"github.com/your-org/vaultdiff/internal/vault"
)

// registerPlanFlags attaches plan-related flags to the given command.
func registerPlanFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("plan", false, "generate an execution plan from the diff")
	cmd.Flags().Bool("plan-noops", false, "include unchanged keys in the plan output")
}

// PlanFlagOptions holds resolved plan flag values.
type PlanFlagOptions struct {
	Enabled      bool
	IncludeNoops bool
}

// resolvePlanOptions reads plan flags from the command and returns options.
func resolvePlanOptions(cmd *cobra.Command) (PlanFlagOptions, vault.PlanOptions) {
	enabled, _ := cmd.Flags().GetBool("plan")
	noops, _ := cmd.Flags().GetBool("plan-noops")
	flagOpts := PlanFlagOptions{
		Enabled:      enabled,
		IncludeNoops: noops,
	}
	planOpts := vault.PlanOptions{
		IncludeNoops: noops,
	}
	return flagOpts, planOpts
}
