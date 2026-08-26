package recovery

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type QuarantineRecord struct {
	Path, Reason string
	MovedAt      time.Time
}

func Quarantine(path, reason string, now time.Time) (QuarantineRecord, error) {
	if path == "" || reason == "" {
		return QuarantineRecord{}, fmt.Errorf("quarantine fields required")
	}
	target := path + fmt.Sprintf(".quarantine.%d", now.UnixNano())
	if err := os.Rename(path, target); err != nil {
		return QuarantineRecord{}, err
	}
	return QuarantineRecord{target, reason, now}, nil
}
func EnsureDir(dir string) error { return os.MkdirAll(filepath.Clean(dir), 0755) }
