package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultdiff/internal/vault"
)

func registerClassifyFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("classify", false, "Enable secret classification")
	cmd.Flags().String("classify-default-tag", "unclassified", "Default tag when no rule matches")
	// Rules are supplied as path=tag or key:prefix=tag pairs, comma-separated.
	cmd.Flags().StringSlice("classify-path-rules", nil, "Path-prefix rules in the form 'prefix=tag' (repeatable)")
	cmd.Flags().StringSlice("classify-key-rules", nil, "Key-prefix rules in the form 'prefix=tag' (repeatable)")
}

func resolveClassifyOptions(cmd *cobra.Command) vault.ClassifyOptions {
	enabled, _ := cmd.Flags().GetBool("classify")
	defaultTag, _ := cmd.Flags().GetString("classify-default-tag")
	pathRules, _ := cmd.Flags().GetStringSlice("classify-path-rules")
	keyRules, _ := cmd.Flags().GetStringSlice("classify-key-rules")

	var rules []vault.ClassifyRule

	for _, r := range pathRules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) == 2 {
			rules = append(rules, vault.ClassifyRule{
				PathPrefix: parts[0],
				Tag:        parts[1],
			})
		}
	}

	for _, r := range keyRules {
		parts := strings.SplitN(r, "=", 2)
		if len(parts) == 2 {
			rules = append(rules, vault.ClassifyRule{
				KeyPrefix: parts[0],
				Tag:       parts[1],
			})
		}
	}

	return vault.ClassifyOptions{
		Enabled:    enabled,
		DefaultTag: defaultTag,
		Rules:      rules,
	}
}
