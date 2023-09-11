package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/calculator"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/handler"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/queue/rabbitmq"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store/memory"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	rabbit, err := rabbitmq.New()
	if err != nil {
		return errors.Wrap(err, "failed rabbit MQ initial")
	}

	m := make(map[string]store.Progression)
	s := memory.New(m)
	calc := calculator.New(rabbit, s)
	cfg := config.Get()
	if cfg.N <= 0 {
		cfg.N = 1
	}

	for i := 0; i < cfg.N; i++ {
		go calc.ConsumeQueue()
	}

	go calc.ClearProgression()

	h := handler.New(calc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/set", h.SetArithmeticProgressionData)
	r.Get("/get", h.GetArithmeticProgressionInfo)
	r.Mount("/debug", middleware.Profiler())

	serve(r)

	return nil
}

func serve(r http.Handler) {
	cfg := config.Get()
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        r,
	}

	err := server.ListenAndServe()
	if err != nil {
		logger.Get().Fatal().Err(err)
	}
}
