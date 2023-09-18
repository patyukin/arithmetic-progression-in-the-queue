package calculator

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/queue"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
	"github.com/pkg/errors"
)

type Params struct {
	N   int     `json:"n"`
	D   float64 `json:"d"`
	Nl  float64 `json:"nl"`
	I   float64 `json:"I"`
	TTL float64 `json:"TTL"`
}

type CalcInterface interface {
	SetProgression(body []byte) error
	ClearProgression()
	ConsumeQueue() error
}

type Calculator struct {
	q queue.Queue
	s store.Store
}

var QueueNumber int

func New(q queue.Queue, s store.Store) CalcInterface {
	return &Calculator{
		q: q,
		s: s,
	}
}

func (c *Calculator) SetProgression(body []byte) error {
	var params Params
	err := json.Unmarshal(body, &params)
	if err != nil {
		return errors.Wrap(err, "failed json.Unmarshal(reqBody, &params) in Calculator.SetProgression")
	}

	// unique uuid
	id := uuid.NewString()
	for c.s.Exists(id) {
		id = uuid.NewString()
	}

	QueueNumber += 1
	c.s.Set(id, store.Progression{
		TTL:              30000,
		Status:           "В очереди",
		QueueNumber:      QueueNumber,
		N:                params.N,
		Nl:               params.Nl,
		I:                params.I,
		D:                params.D,
		CurrentIteration: 0,
		TaskSetUpTime:    time.Now(),
	})

	err = c.q.Publish([]byte(id))

	if err != nil {
		return errors.Wrap(err, "failed Publish in Calculator.SetProgression")
	}

	logger.Get().Info().Msg("Publish message")

	return nil
}

func (c *Calculator) ConsumeQueue() error {
	msgs, err := c.q.Consume()
	if err != nil {
		return errors.Wrap(err, "failed ConsumeQueue")
	}

	var forever chan struct{}

	logger.Get().Info().Msg(" [*] Waiting for messages. To exit press CTRL+C")
	for i := 0; i < config.Get().N; i++ {
		go func() {
			for d := range msgs {
				logger.Get().Info().Msgf("Received a message: %s", d.Body)
				progression, _ := c.s.Get(string(d.Body))
				progression.Status = "В работе"
				c.s.Set(string(d.Body), progression)
				c.s.Loop(string(d.Body))
			}
		}()
	}

	<-forever

	return nil
}

func (c *Calculator) ClearProgression() {
	c.s.ClearTTL()
}
