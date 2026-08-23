package handler

import (
	"github.com/nipalab/nipa/internal/http"
	"github.com/nipalab/nipa/internal/http/model"
)

func (h *Handler) AuthLogin(appCtx http.AppContext) {
	body := &model.LoginRequest{}
	if err := appCtx.ReadJson(body); err != nil {
		appCtx.HandleError(err)
		return
	}
	result, err := h.useCase.AuthUsecase().LoginWithEmailPassword(appCtx.Context(), body.Email, body.Password)
	if err != nil {
		appCtx.HandleError(err)
		return
	}
	appCtx.WriteJson(200, model.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
	})
}

func (h *Handler) AuthRefreshToken(appCtx http.AppContext) {
	body := &model.RefreshTokenRequest{}
	if err := appCtx.ReadJson(body); err != nil {
		appCtx.HandleError(err)
		return
	}
	result, err := h.useCase.AuthUsecase().LoginWithRefreshToken(appCtx.Context(), body.RefreshToken)
	if err != nil {
		appCtx.HandleError(err)
		return
	}
	appCtx.WriteJson(200, model.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		TokenType:    result.TokenType,
	})
}
