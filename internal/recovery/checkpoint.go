package recovery

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type Checkpoint struct {
	EventNumber int64
	Digest      string
}

func Digest(data []byte) string { h := sha256.Sum256(data); return hex.EncodeToString(h[:]) }
func (c Checkpoint) Verify(data []byte) error {
	if c.EventNumber < 0 || c.Digest == "" {
		return fmt.Errorf("invalid checkpoint")
	}
	if Digest(data) != c.Digest {
		return fmt.Errorf("checkpoint digest mismatch")
	}
	return nil
}
