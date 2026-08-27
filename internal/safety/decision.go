package safety

import "time"

type ReleaseDecision struct {
	BlockadeID, Dispatcher string
	Checks                 []string
	Approved               bool
	At                     time.Time
}

func (d ReleaseDecision) Complete() bool {
	return d.BlockadeID != "" && d.Dispatcher != "" && len(d.Checks) >= 3 && d.Approved
}
func (d *ReleaseDecision) Approve() { d.Approved = true; d.At = time.Now() }
func (d *ReleaseDecision) Reject()  { d.Approved = false; d.At = time.Now() }
