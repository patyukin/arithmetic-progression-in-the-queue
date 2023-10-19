package handler

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) GetArithmeticProgressionInfo(w http.ResponseWriter, _ *http.Request) {
	p, err := h.c.GetProgression()
	if err != nil {
		return
	}

	res, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Ошибка при конвертации в JSON:", err)
		return
	}

	_, err = w.Write(res)
	if err != nil {
		return
	}
}
