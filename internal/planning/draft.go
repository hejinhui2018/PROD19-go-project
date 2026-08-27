package planning

import (
	"fmt"
	"github.com/wogo-prod19/railguard/internal/domain"
	"sort"
	"time"
)

type Draft struct {
	ID, BlockadeID, Author string
	Revision               int64
	Segments               []domain.TrackSegment
	Dependencies           []domain.SegmentLink
	Crew                   CrewPlan
	Equipment              EquipmentPlan
	CreatedAt, UpdatedAt   time.Time
}
type CrewPlan struct {
	CrewID    string
	Qualified bool
	Shift     string
	Briefed   bool
	Contact   string
}
type EquipmentPlan struct {
	Required  []string
	Available []string
	Inspected bool
	Custodian string
}
type DependencyIssue struct{ From, To, Reason string }

func NewDraft(id, blockade, author string, segments []domain.TrackSegment, now time.Time) Draft {
	cp := append([]domain.TrackSegment(nil), segments...)
	return Draft{ID: id, BlockadeID: blockade, Author: author, Revision: 1, Segments: cp, CreatedAt: now, UpdatedAt: now}
}
func (d Draft) Validate() error {
	if d.ID == "" || d.BlockadeID == "" || len(d.Segments) == 0 {
		return fmt.Errorf("draft identity and segments required")
	}
	if d.Revision < 1 {
		return fmt.Errorf("revision must be positive")
	}
	return nil
}
func (d Draft) OrderedSegments() []domain.TrackSegment {
	out := append([]domain.TrackSegment(nil), d.Segments...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}
func (d Draft) Ready() bool {
	return d.Crew.Qualified && d.Crew.Briefed && d.Equipment.Inspected && len(d.Equipment.Required) <= len(d.Equipment.Available)
}
