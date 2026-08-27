package planning

import "time"

type Conflict struct {
	DraftID, OtherID, Reason string
	OverlapStart, OverlapEnd time.Time
}

func Overlap(a, b WorkWindow) bool { return a.Start.Before(b.End) && b.Start.Before(a.End) }
func FindConflicts(id string, w WorkWindow, others map[string]WorkWindow) []Conflict {
	out := []Conflict{}
	for oid, ow := range others {
		if oid == id || !Overlap(w, ow) {
			continue
		}
		start := w.Start
		if ow.Start.After(start) {
			start = ow.Start
		}
		end := w.End
		if ow.End.Before(end) {
			end = ow.End
		}
		out = append(out, Conflict{id, oid, "overlapping possession", start, end})
	}
	return out
}
