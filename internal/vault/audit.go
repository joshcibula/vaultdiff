package vault

import (
	"fmt"
	"io"
	"os"
	"time"
)

// AuditOptions controls audit logging of secret access.
type AuditOptions struct {
	Enabled  bool
	Writer   io.Writer
	RedactValues bool
}

// DefaultAuditOptions returns audit logging disabled by default.
func DefaultAuditOptions() AuditOptions {
	return AuditOptions{
		Enabled:      false,
		Writer:       os.Stderr,
		RedactValues: true,
	}
}

// AuditEntry represents a single audited secret access event.
type AuditEntry struct {
	Timestamp time.Time
	Path      string
	Keys      []string
	Source    string
}

// AuditLogger writes audit entries to the configured writer.
type AuditLogger struct {
	opts AuditOptions
}

// NewAuditLogger creates a new AuditLogger with the given options.
func NewAuditLogger(opts AuditOptions) *AuditLogger {
	return &AuditLogger{opts: opts}
}

// Log writes an audit entry if audit logging is enabled.
func (a *AuditLogger) Log(entry AuditEntry) {
	if !a.opts.Enabled {
		return
	}
	keys := entry.Keys
	if a.opts.RedactValues {
		keys = redactKeyNames(keys)
	}
	fmt.Fprintf(
		a.opts.Writer,
		"[audit] %s source=%s path=%s keys=%v\n",
		entry.Timestamp.UTC().Format(time.RFC3339),
		entry.Source,
		entry.Path,
		keys,
	)
}

func redactKeyNames(keys []string) []string {
	out := make([]string, len(keys))
	copy(out, keys)
	return out
}
