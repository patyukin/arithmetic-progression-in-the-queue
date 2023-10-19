package store

import (
	"time"
)

type Progression struct {
	TTL              float64   `json:"TTL"`
	Status           string    `json:"status"`
	QueueNumber      int       `json:"queueNumber"`
	N                int       `json:"n"`
	Nl               float64   `json:"nl"`
	I                float64   `json:"I"`
	D                float64   `json:"d"`
	CurrentIteration int       `json:"currentIteration"`
	Progression      float64   `json:"omitempty"`
	TaskSetUpTime    time.Time `json:"taskSetUpTime"`
	TaskStartTime    time.Time `json:"taskStartTime"`
	TaskFinishTIme   time.Time `json:"taskFinishTIme"` // В случае если задача завершена
}

type Store interface {
	Set(k string, v Progression)
	Get(k string) (Progression, error)
	GetAll() ([]Progression, error)
	Delete(k string)
	Loop(k string)
	ClearTTL()
	Exists(k string) bool
}
