package planning

import "fmt"

type Risk struct {
	Code, Description string
	Severity          int
	Mitigation        string
	Accepted          bool
}
type RiskRegister struct{ Items []Risk }

func (r *RiskRegister) Add(x Risk) error {
	if x.Code == "" || x.Description == "" || x.Severity < 1 {
		return fmt.Errorf("invalid risk")
	}
	r.Items = append(r.Items, x)
	return nil
}
func (r RiskRegister) Open() []Risk {
	out := []Risk{}
	for _, x := range r.Items {
		if !x.Accepted {
			out = append(out, x)
		}
	}
	return out
}
func (r *RiskRegister) Accept(code string) bool {
	for i := range r.Items {
		if r.Items[i].Code == code {
			r.Items[i].Accepted = true
			return true
		}
	}
	return false
}
func (r RiskRegister) Score() int {
	s := 0
	for _, x := range r.Open() {
		s += x.Severity
	}
	return s
}
func (r RiskRegister) Ready() bool { return len(r.Open()) == 0 }
