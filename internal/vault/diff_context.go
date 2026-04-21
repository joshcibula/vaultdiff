package vault

// DiffContext holds metadata about the two sides being compared,
// useful for annotating diff output with source information.
type DiffContext struct {
	LeftPath      string
	RightPath     string
	LeftNamespace string
	RightNamespace string
	LeftMount     string
	RightMount    string
}

// NewDiffContext builds a DiffContext from two parsed VaultPath values.
func NewDiffContext(left, right VaultPath) DiffContext {
	return DiffContext{
		LeftPath:       left.SecretPath,
		RightPath:      right.SecretPath,
		LeftNamespace:  left.Namespace,
		RightNamespace: right.Namespace,
		LeftMount:      left.Mount,
		RightMount:     right.Mount,
	}
}

// SameNamespace returns true when both sides share the same Vault namespace.
func (d DiffContext) SameNamespace() bool {
	return d.LeftNamespace == d.RightNamespace
}

// SameMount returns true when both sides share the same KV mount.
func (d DiffContext) SameMount() bool {
	return d.LeftMount == d.RightMount
}

// Summary returns a human-readable one-line description of the comparison.
func (d DiffContext) Summary() string {
	left := d.LeftPath
	if d.LeftNamespace != "" {
		left = d.LeftNamespace + "/" + left
	}
	right := d.RightPath
	if d.RightNamespace != "" {
		right = d.RightNamespace + "/" + right
	}
	return left + " → " + right
}
