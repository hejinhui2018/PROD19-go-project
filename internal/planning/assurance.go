package planning

import (
	"fmt"
	"time"
)

type Assurance struct {
	DraftID    string
	Required   []string
	Evidence   map[string]string
	Reviewer   string
	ReviewedAt time.Time
}

func (a *Assurance) Add(key, value string) {
	if a.Evidence == nil {
		a.Evidence = map[string]string{}
	}
	a.Evidence[key] = value
}
func (a Assurance) Missing() []string {
	out := []string{}
	for _, k := range a.Required {
		if a.Evidence[k] == "" {
			out = append(out, k)
		}
	}
	return out
}
func (a Assurance) Complete() bool {
	return a.DraftID != "" && a.Reviewer != "" && !a.ReviewedAt.IsZero() && len(a.Missing()) == 0
}
func (a Assurance) Validate() error {
	if !a.Complete() {
		return fmt.Errorf("planning assurance incomplete")
	}
	return nil
}

type Handoff struct {
	From, To string
	Items    []string
	Signed   bool
	At       time.Time
}

func (h *Handoff) Sign(from, to string) bool {
	if from == "" || to == "" || len(h.Items) == 0 {
		return false
	}
	h.From = from
	h.To = to
	h.Signed = true
	h.At = time.Now()
	return true
}
func (h Handoff) Ready() bool { return h.Signed && len(h.Items) > 0 }
