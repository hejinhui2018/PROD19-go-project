package audit

import "strings"

func CSV(records []Record) string {
	rows := []string{"id,actor,action,blockade,from,to,at"}
	for _, r := range records {
		rows = append(rows, strings.Join([]string{r.ID, r.Actor, r.Action, r.BlockadeID, r.From, r.To, r.At.UTC().Format("2006-01-02T15:04:05Z")}, ","))
	}
	return strings.Join(rows, "\n")
}
