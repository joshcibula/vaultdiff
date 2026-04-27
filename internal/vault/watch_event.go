package vault

import "time"

// WatchEventKind describes what changed during a watch cycle.
type WatchEventKind string

const (
	WatchEventAdded    WatchEventKind = "added"
	WatchEventRemoved  WatchEventKind = "removed"
	WatchEventModified WatchEventKind = "modified"
)

// WatchEvent captures a single detected change at a secret path.
type WatchEvent struct {
	Path      string
	Key       string
	Kind      WatchEventKind
	OldValue  string
	NewValue  string
	DetectedAt time.Time
}

// DiffToEvents converts before/after secret maps for a single path into a
// slice of WatchEvents describing what changed.
func DiffToEvents(path string, prev, curr map[string]string) []WatchEvent {
	now := time.Now().UTC()
	var events []WatchEvent

	for k, newVal := range curr {
		oldVal, existed := prev[k]
		if !existed {
			events = append(events, WatchEvent{Path: path, Key: k, Kind: WatchEventAdded, NewValue: newVal, DetectedAt: now})
		} else if oldVal != newVal {
			events = append(events, WatchEvent{Path: path, Key: k, Kind: WatchEventModified, OldValue: oldVal, NewValue: newVal, DetectedAt: now})
		}
	}

	for k, oldVal := range prev {
		if _, stillPresent := curr[k]; !stillPresent {
			events = append(events, WatchEvent{Path: path, Key: k, Kind: WatchEventRemoved, OldValue: oldVal, DetectedAt: now})
		}
	}
	return events
}
