package recovery

import (
	"fmt"
	"time"
)

type RestartReport struct {
	StartedAt, CompletedAt time.Time
	Replay                 ReplayReport
	Compacted              bool
	Quarantined            []QuarantineRecord
	Warnings               []string
}

func (r RestartReport) Summary() string {
	return fmt.Sprintf("recovery events=%d blockades=%d compacted=%t quarantined=%d", r.Replay.Events, r.Replay.Blockades, r.Compacted, len(r.Quarantined))
}
