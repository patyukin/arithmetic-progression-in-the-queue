package store

import (
	"time"
)

type Progression struct {
	TTL              float64
	Status           string
	QueueNumber      int32
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
	Set(k int32, s Progression) error
	GetOneForQueue() (Progression, error)
	GetAll() ([]Progression, error)
}
