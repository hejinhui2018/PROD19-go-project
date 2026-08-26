package recovery

type Metrics struct{ Replays, Corruptions, Compactions, Quarantines int }

func (m *Metrics) RecordReplay(ok bool) {
	m.Replays++
	if !ok {
		m.Corruptions++
	}
}
func (m *Metrics) RecordCompaction() { m.Compactions++ }
func (m *Metrics) RecordQuarantine() { m.Quarantines++ }
