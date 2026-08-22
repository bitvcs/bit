package handler

import (
	"github.com/bitvcs/bit/internal/http"
	"github.com/bitvcs/bit/internal/http/model"
)

func (h *Handler) AuthLogin(appCtx http.AppContext) {
	body := &model.LoginRequest{}
	if err := appCtx.ReadJson(body); err != nil {
		appCtx.WriteJson(400, model.NewAPIError(err.Error()))
		return
	}
	result, err := h.useCase.AuthUsecase().LoginWithEmailPassword(appCtx.Context(), body.Email, body.Password)
	if err != nil {
		appCtx.WriteJson(400, model.NewAPIError(err.Error()))
		return
	}
	appCtx.WriteJson(200, result)
}
