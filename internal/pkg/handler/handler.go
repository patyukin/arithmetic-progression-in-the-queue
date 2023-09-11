package handler

import (
	"io"
	"net/http"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/calculator"
)

type Interface interface {
	SetArithmeticProgressionData(w http.ResponseWriter, r *http.Request)
	GetArithmeticProgressionInfo(w http.ResponseWriter, r *http.Request)
}

type Handler struct {
	c calculator.CalcInterface
}

func New(c calculator.CalcInterface) *Handler {
	return &Handler{
		c: c,
	}
}

func (h *Handler) SetArithmeticProgressionData(_ http.ResponseWriter, r *http.Request) {
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = h.c.SetProgression(reqBody)
	if err != nil {
		return
	}
}

func (h *Handler) GetArithmeticProgressionInfo(_ http.ResponseWriter, _ *http.Request) {
	//h.c.GetProgression()
}
