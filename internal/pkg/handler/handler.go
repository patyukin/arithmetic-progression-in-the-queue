package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/patyukin/arithmetic-progression-in-the-queue/internal/pkg/calculator"
)

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

	var params calculator.Params
	err = json.Unmarshal(reqBody, &params)
	if err != nil {
		return
	}

	err = h.c.SetProgression(params)
	if err != nil {
		return
	}
}

func (h *Handler) GetArithmeticProgressionInfo(_ http.ResponseWriter, _ *http.Request) {
	//h.c.GetProgression()
}
