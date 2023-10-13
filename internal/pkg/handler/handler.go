package handler

import (
	"encoding/json"
	"fmt"
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

func (h *Handler) GetArithmeticProgressionInfo(w http.ResponseWriter, _ *http.Request) {
	p, err := h.c.GetProgression()
	if err != nil {
		fmt.Printf("Ошибка с получением вписка последовательности: %v\n", err)
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
