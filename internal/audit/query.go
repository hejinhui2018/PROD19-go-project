package audit

import "time"

func ByActor(records []Record, actor string) []Record {
	out := []Record{}
	for _, r := range records {
		if r.Actor == actor {
			out = append(out, r)
		}
	}
	return out
}
func Since(records []Record, at time.Time) []Record {
	out := []Record{}
	for _, r := range records {
		if !r.At.Before(at) {
			out = append(out, r)
		}
	}
	return out
}
func LastTransition(records []Record, id string) Record {
	var out Record
	for _, r := range records {
		if r.BlockadeID == id && r.Action == "transition" && r.At.After(out.At) {
			out = r
		}
	}
	return out
}
