package notice

import (
	"fmt"
	"time"
)

type Channel interface{ Send(string, string) error }
type Delivery struct {
	NoticeID string
	Audience Audience
	Attempts int
	SentAt   time.Time
	Error    string
}
type MemoryChannel struct{ Messages []string }

func (c *MemoryChannel) Send(address, message string) error {
	if address == "" || message == "" {
		return fmt.Errorf("empty delivery")
	}
	c.Messages = append(c.Messages, address+":"+message)
	return nil
}
func (d Delivery) Successful() bool { return d.Error == "" && d.Attempts > 0 }
