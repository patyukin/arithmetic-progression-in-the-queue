package store

import (
	"time"
)

type Progression struct {
	TTL              float64
	Status           string
	QueueNumber      int
	N                int
	Nl               float64
	I                float64
	D                float64
	CurrentIteration int
	Progression      float64
	TaskSetUpTime    time.Time
	TaskStartTime    time.Time
	TaskFinishTIme   time.Time // В случае если задача завершена
	Ack              bool
}

type Store interface {
	Add(v Progression) error
	Get(k int) (Progression, error)
	Delete(k int) error
	Loop(k string)
	ClearTTL()
	Exists(k int) bool
}
