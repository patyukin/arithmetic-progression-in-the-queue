package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
)

var _ store.Store = &Memory{}

type Memory struct {
	mx  sync.Mutex
	seq map[string]store.Progression
}

func New(seq map[string]store.Progression) *Memory {
	return &Memory{
		seq: seq,
	}
}

func (m *Memory) Exists(k string) bool {
	m.mx.Lock()
	defer m.mx.Unlock()
	mp, ok := m.seq[k]
	if !ok {
		return false
	}

	if mp.Status == "Завершена" {
		now := time.Now()
		finishTime := mp.TaskSetUpTime.Add(time.Duration(int(mp.TTL)) * time.Second)
		if finishTime.Sub(now) > 0 {
			delete(m.seq, k)

			return false
		}
	}

	return true
}

func (m *Memory) Set(k string, v store.Progression) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	if v.Status == "Завершена" {
		now := time.Now()
		finishTime := v.TaskSetUpTime.Add(time.Duration(int(v.TTL)) * time.Second)

		if finishTime.Sub(now) > 0 {
			delete(m.seq, k)
		}

		return nil
	}

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

	if m.seq[k].Status == "Завершена" {
		now := time.Now()
		finishTime := v.TaskSetUpTime.Add(time.Duration(int(v.TTL)) * time.Second)

		if finishTime.Sub(now) > 0 {
			delete(m.seq, k)
		}
	}

	return v, nil
}

func (m *Memory) Delete(k string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.seq, k)
}

func (m *Memory) ClearTTL() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
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
}

func (m *Memory) Loop(k string) {
	m.mx.Lock()
	current := m.seq[k]
	m.mx.Unlock()

	logger.Get().Info().Msgf("Start map with key: %s", k)
	current.Progression = current.Nl
	for i := 0; i < current.N; i++ {
		current.CurrentIteration++
		current.Progression += current.D
		m.seq[k] = current

		time.Sleep(time.Duration(current.I) * time.Second)
	}

	current.Status = "Завершен"

	m.mx.Lock()
	m.seq[k] = current
	m.mx.Unlock()
	logger.Get().Info().Msgf("Finish map with key: %s", k)
}
