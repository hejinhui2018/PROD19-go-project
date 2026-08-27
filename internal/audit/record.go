package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Record struct {
	ID, Actor, Action, BlockadeID, From, To, Summary, Hash string
	At                                                     time.Time
}

func New(id, actor, action, bid, from, to, summary string, at time.Time) Record {
	r := Record{id, actor, action, bid, from, to, summary, "", at}
	h := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d", id, actor, action, bid, from, to, at.UnixNano())))
	r.Hash = hex.EncodeToString(h[:])
	return r
}
