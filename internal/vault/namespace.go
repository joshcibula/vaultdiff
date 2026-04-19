package vault

import (
	"fmt"
	"strings"
)

// ParsePath splits a vault path into mount, namespace, and secret path components.
// Accepted formats:
//   - secret/data/myapp
//   - ns1/secret/data/myapp  (when namespace prefix is detected)
type ParsedPath struct {
	Namespace string
	Mount     string
	SecretPath string
}

// ParseVaultPath parses a raw vault path string into its components.
// If a namespace is provided explicitly it takes precedence.
func ParseVaultPath(raw, defaultNamespace string) (*ParsedPath, error) {
	raw = strings.Trim(raw, "/")
	if raw == "" {
		return nil, fmt.Errorf("vault path must not be empty")
	}

	parts := strings.SplitN(raw, "/", 3)
	if len(parts) < 2 {
		return nil, fmt.Errorf("vault path %q is too short: expected at least mount/secret-path", raw)
	}

	p := &ParsedPath{
		Namespace: defaultNamespace,
	}

	// Heuristic: if the path has 3+ segments and the first segment contains no
	// "data" keyword, treat the first segment as a namespace override.
	if len(parts) == 3 && parts[1] != "data" && parts[1] != "metadata" {
		p.Namespace = parts[0]
		p.Mount = parts[1]
		p.SecretPath = parts[2]
	} else {
		p.Mount = parts[0]
		p.SecretPath = strings.Join(parts[1:], "/")
	}

	return p, nil
}

// FullKVv2Path returns the KV v2 data path for use with the Vault API.
func (p *ParsedPath) FullKVv2Path() string {
	// Ensure the mount path includes the /data/ prefix required by KV v2.
	if strings.HasPrefix(p.SecretPath, "data/") {
		return fmt.Sprintf("%s/%s", p.Mount, p.SecretPath)
	}
	return fmt.Sprintf("%s/data/%s", p.Mount, p.SecretPath)
}
