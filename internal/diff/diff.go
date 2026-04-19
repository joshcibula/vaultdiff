package diff

import "fmt"

// SecretMap represents a flat map of secret key-value pairs at a path.
type SecretMap map[string]string

// Result holds the diff outcome between two secret maps.
type Result struct {
	OnlyInLeft  map[string]string
	OnlyInRight map[string]string
	Modified    map[string][2]string // key -> [leftVal, rightVal]
	Unchanged   map[string]string
}

// Compare diffs two SecretMaps and returns a Result.
func Compare(left, right SecretMap) Result {
	res := Result{
		OnlyInLeft:  make(map[string]string),
		OnlyInRight: make(map[string]string),
		Modified:    make(map[string][2]string),
		Unchanged:   make(map[string]string),
	}

	for k, lv := range left {
		if rv, ok := right[k]; ok {
			if lv == rv {
				res.Unchanged[k] = lv
			} else {
				res.Modified[k] = [2]string{lv, rv}
			}
		} else {
			res.OnlyInLeft[k] = lv
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			res.OnlyInRight[k] = rv
		}
	}

	return res
}

// HasDifferences returns true if there are any differences.
func (r Result) HasDifferences() bool {
	return len(r.OnlyInLeft) > 0 || len(r.OnlyInRight) > 0 || len(r.Modified) > 0
}

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	return fmt.Sprintf(
		"added: %d, removed: %d, modified: %d, unchanged: %d",
		len(r.OnlyInRight),
		len(r.OnlyInLeft),
		len(r.Modified),
		len(r.Unchanged),
	)
}
