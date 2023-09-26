package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
)

var _ store.Store = &Memory{}

type Memory struct {
	mx  sync.Mutex
	s   []store.Progression
	cfg *config.Config
	l   *logger.Logger
}

func New(s []store.Progression, cfg *config.Config, l *logger.Logger) *Memory {
	return &Memory{
		s:   s,
		cfg: cfg,
		l:   l,
	}
}

func (m *Memory) Exists(k int) bool {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, v := range m.s {
		if k == v.QueueNumber {
			return true
		}
	}

	return false
}

func (m *Memory) Add(v store.Progression) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.s = append(m.s, v)
	return nil
}

func (m *Memory) Get(k int) (store.Progression, error) {
	m.mx.Lock()
	defer m.mx.Unlock()

	for _, v := range m.s {
		if k == v.QueueNumber {
			return v, nil
		}
	}

	return store.Progression{}, fmt.Errorf("failed to find key: %d", k)
}

func (m *Memory) Delete(k int) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	var idx int
	for i, v := range m.s {
		if k == v.QueueNumber {
			idx = i
		}
	}

	m.s = append(m.s[:idx], m.s[idx+1:]...)

	return nil
}

func (m *Memory) ClearTTL() {
	for {
		m.mx.Lock()
		for k, v := range m.s {
			if v.Status != "Завершена" {
				continue
			}

			now := time.Now()
			finishTime := v.TaskSetUpTime.Add(time.Duration(int(v.TTL)) * time.Second)

			if finishTime.Sub(now) > 0 {
				delete(m.s, k)
			}
		}

		m.mx.Unlock()
	}
}

func (m *Memory) Loop(k string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	current := m.s[k]
	current.Progression = current.Nl
	for i := 0; i < current.N; i++ {
		current.CurrentIteration++
		current.Progression += current.D
		m.s[k] = current
		time.Sleep(time.Duration(current.I) * time.Second)
	}

	current.Status = "Завершен"
	m.s[k] = current
}
