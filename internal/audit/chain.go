package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Chain(records []Record) string {
	seed := ""
	for _, r := range records {
		h := sha256.Sum256([]byte(seed + "|" + r.Hash))
		seed = hex.EncodeToString(h[:])
	}
	return seed
}
func VerifyChain(records []Record, expected string) error {
	if Chain(records) != expected {
		return fmt.Errorf("audit chain mismatch")
	}
	return nil
}
