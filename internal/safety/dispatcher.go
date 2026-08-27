package safety

import "fmt"

type Hold struct {
	ID, Reason, Dispatcher string
	Active                 bool
}

func PlaceHold(id, reason, dispatcher string) (Hold, error) {
	if id == "" || reason == "" || dispatcher == "" {
		return Hold{}, fmt.Errorf("hold fields required")
	}
	return Hold{id, reason, dispatcher, true}, nil
}
func ReleaseHold(h Hold, dispatcher string) error {
	if !h.Active {
		return fmt.Errorf("hold already released")
	}
	if h.Dispatcher != dispatcher {
		return fmt.Errorf("dispatcher attribution mismatch")
	}
	return nil
}
