package api

import (
	"context"

	"github.com/bitvcs/bit/internal/domain"
	"github.com/bitvcs/bit/internal/http"
	"github.com/emicklei/go-restful/v3"
)

type appContext struct {
	req    *restful.Request
	resp   *restful.Response
	claims *domain.Claims
}

func newAppContext(req *restful.Request, resp *restful.Response) (http.AppContext, error) {
	var claims *domain.Claims
	claimInterface, ok := req.Attribute(AttributeClaims).(*domain.Claims)
	if ok {
		claims = claimInterface
	}
	return &appContext{
		req:    req,
		resp:   resp,
		claims: claims,
	}, nil
}

func (a *appContext) Context() context.Context {
	return a.req.Request.Context()
}

func (a *appContext) Claims() *domain.Claims {
	return a.claims
}

func (a *appContext) ReadJson(v any) error {
	return a.req.ReadEntity(v)
}

func (a *appContext) WriteJson(statusCode int, v any) error {
	return a.resp.WriteHeaderAndJson(statusCode, v, restful.MIME_JSON)
}
