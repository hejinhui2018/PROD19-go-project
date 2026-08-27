package safety

type Role string

const (
	RoleDispatcher Role = "dispatcher"
	RoleSignal     Role = "signal"
	RoleCrew       Role = "crew"
	RoleAuthority  Role = "authority"
)

type Roster struct{ Members map[Role][]string }

func (r Roster) Has(role Role, name string) bool {
	for _, x := range r.Members[role] {
		if x == name {
			return true
		}
	}
	return false
}
func (r *Roster) Add(role Role, name string) {
	if r.Members == nil {
		r.Members = map[Role][]string{}
	}
	if !r.Has(role, name) {
		r.Members[role] = append(r.Members[role], name)
	}
}
