package vault

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"
)

// ExportFormat defines the output format for exported secrets.
type ExportFormat string

const (
	ExportFormatJSON ExportFormat = "json"
	ExportFormatCSV  ExportFormat = "csv"
	ExportFormatEnv  ExportFormat = "env"
)

// ExportOptions configures how secrets are exported.
type ExportOptions struct {
	Format    ExportFormat
	PathLabel string // column/field name for the path key
}

// DefaultExportOptions returns sensible defaults.
func DefaultExportOptions() ExportOptions {
	return ExportOptions{
		Format:    ExportFormatJSON,
		PathLabel: "path",
	}
}

// ExportSecrets writes the provided secrets map to w in the configured format.
// secrets is a map of path -> (key -> value).
func ExportSecrets(w io.Writer, secrets map[string]map[string]string, opts ExportOptions) error {
	switch opts.Format {
	case ExportFormatJSON:
		return exportJSON(w, secrets)
	case ExportFormatCSV:
		return exportCSV(w, secrets, opts.PathLabel)
	case ExportFormatEnv:
		return exportEnv(w, secrets)
	default:
		return fmt.Errorf("unsupported export format: %q", opts.Format)
	}
}

func exportJSON(w io.Writer, secrets map[string]map[string]string) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func exportCSV(w io.Writer, secrets map[string]map[string]string, pathLabel string) error {
	cw := csv.NewWriter(w)
	paths := sortedKeys(secrets)
	for _, path := range paths {
		kv := secrets[path]
		keys := make([]string, 0, len(kv))
		for k := range kv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if err := cw.Write([]string{path, k, kv[k]}); err != nil {
				return err
			}
		}
	}
	cw.Flush()
	return cw.Error()
}

func exportEnv(w io.Writer, secrets map[string]map[string]string) error {
	paths := sortedKeys(secrets)
	for _, path := range paths {
		kv := secrets[path]
		keys := make([]string, 0, len(kv))
		for k := range kv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if _, err := fmt.Fprintf(w, "%s=%s\n", k, kv[k]); err != nil {
				return err
			}
		}
	}
	return nil
}

func sortedKeys(m map[string]map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
