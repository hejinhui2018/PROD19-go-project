package ops

import (
	"fmt"
	"sync/atomic"
	"time"
)

var seq uint64

func ID(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), atomic.AddUint64(&seq, 1))
}
