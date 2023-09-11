package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
)

var _ store.Store = &Memory{}

type Memory struct {
	mx  sync.Mutex
	seq map[string]store.Progression
}

func (m *Memory) SetByParams(k string, v store.Progression) error {
	//TODO implement me
	panic("implement me")
}

func New(seq map[string]store.Progression) *Memory {
	return &Memory{
		seq: seq,
	}
}

func (m *Memory) Exists(k string) bool {
	m.mx.Lock()
	defer m.mx.Unlock()

	if _, ok := m.seq[k]; ok {
		return true
	}

	return false
}

func (m *Memory) Set(k string, v store.Progression) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.seq[k] = v
	return nil
}

func (m *Memory) Get(k string) (store.Progression, error) {
	m.mx.Lock()
	defer m.mx.Unlock()
	v, ok := m.seq[k]
	if !ok {
		return store.Progression{}, fmt.Errorf("failed to find key: %s", k)
	}

	return v, nil
}

func (m *Memory) Delete(k string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.seq, k)
}

func (m *Memory) ClearTTL() {
	for {
		m.mx.Lock()
		for k, v := range m.seq {
			if v.Status != "Завершена" {
				continue
			}

			now := time.Now()
			finishTime := v.TaskSetUpTime.Add(time.Duration(int(v.TTL)) * time.Second)

			if finishTime.Sub(now) > 0 {
				delete(m.seq, k)
			}
		}

		m.mx.Unlock()
	}
}

func (m *Memory) Loop(k string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	current := m.seq[k]
	current.Progression = current.Nl
	for i := 0; i < current.N; i++ {
		current.CurrentIteration++
		current.Progression += current.D
		m.seq[k] = current
		time.Sleep(time.Duration(current.I) * time.Second)
	}

	current.Status = "Завершен"
	m.seq[k] = current
}
