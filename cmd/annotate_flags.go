package cmd

import (
	"github.com/spf13/cobra"
	"github.com/wndhydrnt/vaultdiff/internal/vault"
)

// registerAnnotateFlags attaches annotation-related flags to a command.
func registerAnnotateFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("annotate", false, "inject source metadata tags into each secret")
	cmd.Flags().String("annotate-tag-key", "_vaultdiff_source", "key name used for the injected source tag")
	cmd.Flags().String("annotate-tag-value", "", "static value for the source tag (defaults to secret path)")
	cmd.Flags().String("annotate-path-prefix", "", "path prefix to strip when deriving tag value from path")
	cmd.Flags().StringToString("annotate-custom-tags", map[string]string{}, "additional key=value tags injected into every secret")
}

// resolveAnnotateOptions builds AnnotateOptions from parsed command flags.
func resolveAnnotateOptions(cmd *cobra.Command) (vault.AnnotateOptions, error) {
	opts := vault.DefaultAnnotateOptions()

	enabled, err := cmd.Flags().GetBool("annotate")
	if err != nil {
		return opts, err
	}
	opts.Enabled = enabled

	tagKey, err := cmd.Flags().GetString("annotate-tag-key")
	if err != nil {
		return opts, err
	}
	opts.TagKey = tagKey

	tagValue, err := cmd.Flags().GetString("annotate-tag-value")
	if err != nil {
		return opts, err
	}
	opts.TagValue = tagValue

	pathPrefix, err := cmd.Flags().GetString("annotate-path-prefix")
	if err != nil {
		return opts, err
	}
	opts.PathPrefix = pathPrefix

	customTags, err := cmd.Flags().GetStringToString("annotate-custom-tags")
	if err != nil {
		return opts, err
	}
	opts.CustomTags = customTags

	return opts, nil
}
