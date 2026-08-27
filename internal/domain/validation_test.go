package domain

import (
	"testing"
	"time"
)

func TestValidate(t *testing.T) {
	b := Blockade{ID: "x", Name: "n", Segment: TrackSegment{ID: "s"}, Crew: Crew{ID: "c", Qualified: true}, Start: time.Now(), End: time.Now().Add(time.Hour)}
	if ValidateBlockade(b) != nil {
		t.Fatal()
	}
}
