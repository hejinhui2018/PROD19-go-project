package domain

import "fmt"

func ValidateBlockade(b Blockade) error {
	if b.ID == "" || b.Name == "" || b.Segment.ID == "" || b.Crew.ID == "" {
		return fmt.Errorf("%w: identity fields required", ErrInvalid)
	}
	if !b.Start.Before(b.End) {
		return fmt.Errorf("%w: start must precede end", ErrInvalid)
	}
	if !b.Crew.Qualified {
		return fmt.Errorf("%w: crew is not qualified", ErrInvalid)
	}
	return nil
}
