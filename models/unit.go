package models

import (
	"time"

	"github.com/AltMax/art-test/services"
)

type Unit struct {
	ID        string
	Data      []byte
	CreatedAt time.Time
}

func (u *Unit) Proto() *services.Unit {
	if u == nil {
		return nil
	}
	return &services.Unit{
		Id:        u.ID,
		Data:      u.Data,
		CreatedAt: timeToMilliseconds(u.CreatedAt),
	}
}

type Units []*Unit

func (us Units) Proto() []*services.Unit {
	pb := make([]*services.Unit, 0, len(us))
	for _, u := range us {
		pb = append(pb, u.Proto())
	}
	return pb
}

func timeToMilliseconds(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
