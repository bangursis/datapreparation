package delivery

import (
	"datapreparation/profiles"
	"datapreparation/profiles/usecase"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/valyala/fasthttp"
)

type fastHTTP struct {
	uc profiles.UseCase
}

type request struct {
	FilePath string `json:"filepath"`
}

func NewHandler(uc profiles.UseCase) *fastHTTP {

	return &fastHTTP{
		uc: uc,
	}
}

func (h *fastHTTP) ImportICCID(ctx *fasthttp.RequestCtx) {
	r := &request{}
	err := json.Unmarshal(ctx.PostBody(), r)
	if err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	err = h.uc.Import(ctx, r.FilePath)
	if err != nil {
		if errors.Is(err, usecase.DecryptionErr) {
			ctx.SetStatusCode(http.StatusBadRequest)
		}
		ctx.SetStatusCode(http.StatusInternalServerError)
		return
	}
}
