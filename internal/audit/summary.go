package audit

type TransitionSummary struct {
	BlockadeID string
	Counts     map[string]int
	Actors     map[string]int
}

func Summarize(records []Record, id string) TransitionSummary {
	s := TransitionSummary{BlockadeID: id, Counts: map[string]int{}, Actors: map[string]int{}}
	for _, r := range records {
		if r.BlockadeID == id {
			s.Counts[r.Action]++
			s.Actors[r.Actor]++
		}
	}
	return s
}
