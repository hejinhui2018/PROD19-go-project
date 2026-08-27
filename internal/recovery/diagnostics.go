package recovery

import "fmt"

type Diagnostic struct {
	Component, Status, Detail string
	Critical                  bool
}
type DiagnosticSet struct{ Items []Diagnostic }

func (d *DiagnosticSet) Add(x Diagnostic) { d.Items = append(d.Items, x) }
func (d DiagnosticSet) Healthy() bool {
	for _, x := range d.Items {
		if x.Critical && x.Status != "ok" {
			return false
		}
	}
	return true
}
func (d DiagnosticSet) Summary() string {
	bad := 0
	for _, x := range d.Items {
		if x.Status != "ok" {
			bad++
		}
	}
	return fmt.Sprintf("diagnostics=%d failed=%d", len(d.Items), bad)
}
