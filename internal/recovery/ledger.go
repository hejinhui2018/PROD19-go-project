package recovery

type RecoveryLedger struct{ Reports []RestartReport }

func (l *RecoveryLedger) Add(r RestartReport) { l.Reports = append(l.Reports, r) }
func (l RecoveryLedger) Latest() RestartReport {
	if len(l.Reports) == 0 {
		return RestartReport{}
	}
	return l.Reports[len(l.Reports)-1]
}
func (l RecoveryLedger) Healthy() bool {
	for _, r := range l.Reports {
		if len(r.Replay.Errors) > 0 {
			return false
		}
	}
	return true
}
