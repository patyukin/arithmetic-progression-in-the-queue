package calculator

import (
	"context"
	"time"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
	"github.com/pkg/errors"
)

const (
	IN_PROCESS = "В процессе"
	IN_QUEUE   = "В очереди"
	COMPLETED  = "Завершена"
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
	ClearProgression()
	ConsumeQueue(ctx context.Context) error
}

type Calculator struct {
	s   store.Store
	cfg *config.Config
	l   *logger.Logger
}

var QueueNumber int

func New(s store.Store, cfg *config.Config, l *logger.Logger) *Calculator {
	return &Calculator{
		s:   s,
		cfg: cfg,
		l:   l,
	}
}

func (c *Calculator) SetProgression(params Params) error {
	QueueNumber += 1
	err := c.s.Add(store.Progression{
		TTL:              30000,
		Status:           IN_QUEUE,
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

func (c *Calculator) ConsumeQueue(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			for i := 0; i < c.cfg.N; i++ {
				go func() {
					progression, _ := c.s.Get(string(d.Body))
					progression.Status = IN_PROCESS
					c.s.Add(string(d.Body), progression)
					go c.s.Loop(string(d.Body))
				}()
			}
		}
	}
}

func (c *Calculator) ClearProgression() {
	c.s.ClearTTL()
}
