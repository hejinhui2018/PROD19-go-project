package domain

import "time"

type State string

const (
	StateDraft          State = "draft"
	StatePlanned        State = "planned"
	StateCrewReady      State = "crew_ready"
	StateProtected      State = "protected"
	StateActive         State = "active"
	StateReleasePending State = "release_pending"
	StateReleased       State = "released"
	StateCancelled      State = "cancelled"
	StateIncident       State = "incident"
)

type TrackSegment struct{ ID, Name string }
type SegmentLink struct {
	From, To string
	Distance int
}
type Crew struct {
	ID, Name  string
	Qualified bool
}
type SafetyLockPack struct {
	ID        string
	Applied   bool
	AppliedAt time.Time
}
type DispatcherApproval struct {
	ID, Dispatcher string
	Approved       bool
	At             time.Time
}
type EquipmentHandoff struct {
	ID, Custodian string
	Complete      bool
	At            time.Time
}
type FieldObservation struct {
	ID, Reporter, Message string
	At                    time.Time
}
type ReleaseNotice struct {
	ID, BlockadeID, Audience, Message string
	Sent                              bool
	CreatedAt                         time.Time
}
type Blockade struct {
	ID, Name     string
	Segment      TrackSegment
	Crew         Crew
	Locks        SafetyLockPack
	Approval     DispatcherApproval
	Equipment    EquipmentHandoff
	State        State
	Start, End   time.Time
	Observations []FieldObservation
	Version      int64
	UpdatedAt    time.Time
	PlanRevision int64
	Adjacent     []SegmentLink
}
type Event struct {
	Number     int64     `json:"number"`
	Type       string    `json:"type"`
	BlockadeID string    `json:"blockade_id"`
	Payload    Blockade  `json:"payload"`
	At         time.Time `json:"at"`
}
