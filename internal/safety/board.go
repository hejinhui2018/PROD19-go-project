package safety

import (
	"fmt"
	"time"
)

type ProtectionBoard struct {
	ID, BlockadeID, IssuedBy string
	Entries                  []BoardEntry
	IssuedAt                 time.Time
	Active                   bool
}
type BoardEntry struct {
	Signal, Aspect, Route string
	Locked                bool
}

func IssueBoard(id, blockade, issuer string, entries []BoardEntry, now time.Time) (ProtectionBoard, error) {
	if id == "" || blockade == "" || issuer == "" || len(entries) == 0 {
		return ProtectionBoard{}, fmt.Errorf("board issuance incomplete")
	}
	for _, e := range entries {
		if e.Signal == "" || !e.Locked {
			return ProtectionBoard{}, fmt.Errorf("unlocked board entry")
		}
	}
	return ProtectionBoard{id, blockade, issuer, entries, now, true}, nil
}
func (b ProtectionBoard) Close() ProtectionBoard { b.Active = false; return b }
