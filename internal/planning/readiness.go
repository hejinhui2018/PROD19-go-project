package planning

import "strings"

func MissingEquipment(p EquipmentPlan) []string {
	have := map[string]bool{}
	for _, x := range p.Available {
		have[strings.ToLower(x)] = true
	}
	out := []string{}
	for _, x := range p.Required {
		if !have[strings.ToLower(x)] {
			out = append(out, x)
		}
	}
	return out
}
func CrewReady(p CrewPlan) bool { return p.CrewID != "" && p.Qualified && p.Briefed && p.Shift != "" }
