package planning

type ChecklistItem struct {
	Key, Label         string
	Required, Complete bool
	Owner              string
}
type Checklist struct{ Items []ChecklistItem }

func (c Checklist) Complete() bool {
	for _, i := range c.Items {
		if i.Required && !i.Complete {
			return false
		}
	}
	return true
}
func (c Checklist) Missing() []ChecklistItem {
	out := []ChecklistItem{}
	for _, i := range c.Items {
		if i.Required && !i.Complete {
			out = append(out, i)
		}
	}
	return out
}
func (c *Checklist) Mark(key, owner string) bool {
	for i := range c.Items {
		if c.Items[i].Key == key && (c.Items[i].Owner == "" || c.Items[i].Owner == owner) {
			c.Items[i].Complete = true
			return true
		}
	}
	return false
}
