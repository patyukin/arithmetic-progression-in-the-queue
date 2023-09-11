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
}

type Store interface {
	Set(k string, v Progression) error
	SetByParams(k string, v Progression) error
	Get(k string) (Progression, error)
	Delete(k string)
	Loop(k string)
	ClearTTL()
	Exists(k string) bool
}
