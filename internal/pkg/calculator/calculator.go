package calculator

import (
	"sync/atomic"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
	"github.com/pkg/errors"
)

const (
	InProcess = "В процессе"
	InQueue   = "В очереди"
	Completed = "Завершена"
)

type Params struct {
	N   int     `json:"n"`
	D   float64 `json:"d"`
	I   float64 `json:"I"`
	Nl  float64 `json:"nl"`
	TTL float64 `json:"TTL"`
}

type CalcInterface interface {
	SetProgression(params Params) error
	ConsumeQueue() error
	GetProgression() ([]store.Progression, error)
}

type Calculator struct {
	s   store.Store
	cfg *config.Config
	l   *logger.Logger
}

var QueueNumber int32

func New(s store.Store, cfg *config.Config, l *logger.Logger) *Calculator {
	return &Calculator{
		s:   s,
		cfg: cfg,
		l:   l,
	}
}

func (c *Calculator) SetProgression(params Params) error {
	atomic.AddInt32(&QueueNumber, 1)
	err := c.s.Add(store.Progression{
		TTL:              30000,
		Status:           InQueue,
		QueueNumber:      QueueNumber,
		N:                params.N,
		Nl:               params.Nl,
		I:                params.I,
		D:                params.D,
		CurrentIteration: 0,
		TaskSetUpTime:    time.Now(),
	})

	if err != nil {
		return errors.Wrap(err, "failed store.Store.SetProgression in Calculator.SetProgression")
	}

	return nil
}

func (c *Calculator) ConsumeQueue() error {
	for i := 0; i < c.cfg.N; i++ {
		go func() {
			current, _ := c.s.GetOneForQueue()

			for i := 0; i < current.N; i++ {
				current.CurrentIteration++
				current.Progression += current.D
				time.Sleep(time.Duration(current.I) * time.Second)
			}

			current.Status = Completed

			err := c.s.Set(current.QueueNumber, current)
			if err != nil {
				c.l.Error().Msgf("filed to set slice: %v", err)
			}
		}()
	}

	return nil
}

func (c *Calculator) GetProgression() ([]store.Progression, error) {
	return c.s.GetAll()
}
