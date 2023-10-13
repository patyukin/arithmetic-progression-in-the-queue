package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/calculator"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
)

var _ store.Store = &Memory{}

type Memory struct {
	firstInQueue int
	mx           sync.RWMutex
	s            []store.Progression
	cfg          *config.Config
	l            *logger.Logger
}

func New(ctx context.Context, s []store.Progression, cfg *config.Config, l *logger.Logger) *Memory {
	m := &Memory{
		s:   s,
		cfg: cfg,
		l:   l,
	}

	go m.clearTTL(ctx)

	return m
}

func (m *Memory) Add(v store.Progression) error {
	m.s = append(m.s, v)
	return nil
}

func (m *Memory) GetOneForQueue() (store.Progression, error) {
	m.mx.Lock()
	if len(m.s) <= m.firstInQueue {
		return store.Progression{}, fmt.Errorf("queue not found")
	}
	m.s[m.firstInQueue].Status = calculator.InProcess
	result := m.s[m.firstInQueue]
	m.firstInQueue++
	defer m.mx.Unlock()

	return result, nil
}

func (m *Memory) Set(k int32, s store.Progression) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	var idx int
	for i, v := range m.s {
		if k == v.QueueNumber {
			idx = i
		}
	}

	m.s[idx] = s

	return nil
}

func (m *Memory) clearTTL(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 10)

	for {
		select {
		case <-ticker.C:
			m.mx.Lock()
			for i := 0; i < len(m.s[:m.firstInQueue]); i++ {
				current := m.s[i]
				if current.Status != calculator.Completed {
					continue
				}

				now := time.Now()
				finishTime := current.TaskSetUpTime.Add(time.Duration(int(current.TTL)) * time.Second)
				if finishTime.Sub(now) > 0 {
					m.s = append(m.s[:i], m.s[i+1:]...)
					i--
				}
			}

			m.mx.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (m *Memory) GetAll() ([]store.Progression, error) {
	m.mx.Lock()
	progressions := make([]store.Progression, 0, len(m.s))
	for _, v := range m.s {
		progressions = append(progressions, store.Progression{
			TTL:              v.TTL,
			Status:           v.Status,
			QueueNumber:      v.QueueNumber,
			N:                v.N,
			Nl:               v.Nl,
			I:                v.I,
			D:                v.D,
			CurrentIteration: v.CurrentIteration,
			Progression:      v.Progression,
			TaskSetUpTime:    v.TaskSetUpTime,
			TaskStartTime:    v.TaskStartTime,
			TaskFinishTIme:   v.TaskFinishTIme,
		})
	}

	m.mx.Unlock()

	sort.SliceStable(progressions, func(i, j int) bool {
		return progressions[i].QueueNumber < progressions[j].QueueNumber
	})

	return progressions, nil
}
