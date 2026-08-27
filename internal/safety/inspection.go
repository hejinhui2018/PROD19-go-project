package safety

import "time"

type Inspection struct {
	ID, Inspector string
	Items         map[string]bool
	At            time.Time
}

func (i Inspection) Passed(required []string) bool {
	for _, k := range required {
		if !i.Items[k] {
			return false
		}
	}
	return true
}
func (i *Inspection) Set(item string, ok bool) {
	if i.Items == nil {
		i.Items = map[string]bool{}
	}
	i.Items[item] = ok
}
func (i Inspection) Failed() []string {
	out := []string{}
	for k, v := range i.Items {
		if !v {
			out = append(out, k)
		}
	}
	return out
}
