package recovery

type CompactionPolicy struct {
	MaxEvents int
	MaxBytes  int64
	KeepTail  int
}

func (p CompactionPolicy) ShouldCompact(events int, bytes int64) bool {
	return (p.MaxEvents > 0 && events >= p.MaxEvents) || (p.MaxBytes > 0 && bytes >= p.MaxBytes)
}
func (p CompactionPolicy) Retained(total int) int {
	if p.KeepTail < 0 {
		return 0
	}
	if total < p.KeepTail {
		return total
	}
	return p.KeepTail
}
