package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/calculator"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/handler"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store"
	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/store/memory"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/config"
	"github.com/patyukin/arithmetic-progression-in-the-queue/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg := config.Init()
	l := logger.Init(cfg)

	var sl []store.Progression
	s := memory.New(ctx, sl, cfg, l)
	calc := calculator.New(s, cfg, l)
	if cfg.N <= 0 {
		cfg.N = 1
	}

	go calc.ConsumeQueue()

	h := handler.New(calc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Post("/set", h.SetArithmeticProgressionData)
	r.Get("/get", h.GetArithmeticProgressionInfo)
	r.Mount("/debug", middleware.Profiler())

	serve(ctx, r, cfg)

	<-ctx.Done()
	closeDeps()

	return nil
}

func serve(ctx context.Context, r http.Handler, cfg *config.Config) {
	go func() {
		addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
		server := &http.Server{
			Addr:           addr,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
			Handler:        r,
		}

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func closeDeps() {

}
