package api

import (
	"net/http"

	httpApp "github.com/bitvcs/bit/internal/http"
	"github.com/bitvcs/bit/internal/http/handler"
	"github.com/bitvcs/bit/internal/http/model"
	"github.com/bitvcs/bit/internal/http/swagger"
	"github.com/bitvcs/bit/internal/usecase"
	"github.com/emicklei/go-restful/v3"
)

type usecaseContainer interface {
	AuthUsecase() *usecase.Auth
	UserUsecase() *usecase.User
}

type API struct {
	useCase usecaseContainer
}

func NewAPI(useCase usecaseContainer) *API {
	return &API{
		useCase: useCase,
	}
}

func (a *API) SetupRoute() http.Handler {
	//authFilter := NewAuthFilter(a.useCase.AuthUsecase())
	cors := restful.CrossOriginResourceSharing{
		AllowedMethods: []string{"POST", "GET", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "Accept"},
	}

	handler := handler.NewHandler(a.useCase)

	authWs := new(restful.WebService).ApiVersion("1.0.0")
	authWs.Path("/auth").
		Filter(cors.Filter).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	setupAuthRouter(authWs, handler)
	restful.Add(authWs)

	swagger.SetupSwagger()

	return restful.DefaultContainer
}

func wrap(fn func(httpApp.AppContext)) restful.RouteFunction {
	return func(req *restful.Request, resp *restful.Response) {
		ctx, err := newAppContext(req, resp)
		if err != nil {
			resp.WriteHeaderAndJson(http.StatusInternalServerError, model.NewAPIError(err.Error()), restful.MIME_JSON)
			return
		}
		fn(ctx)
	}
}
