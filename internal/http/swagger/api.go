package swagger

import (
	"embed"
	"io/fs"
	"net/http"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

//go:embed swagger-ui/*
var swaggerUi embed.FS

func SetupSwagger() {
	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/docs/api.json",
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       "Bitd Docs",
					Description: "Bitd Docs",
				},
			}
			s.SecurityDefinitions = map[string]*spec.SecurityScheme{
				"authorization": spec.APIKeyAuth("Authorization", "header"),
			}
			s.Security = []map[string][]string{
				{"authorization": {}},
			}
		},
	}
	fs, _ := fs.Sub(swaggerUi, "swagger-ui")
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.FS(fs))))
	service := restfulspec.NewOpenAPIService(config)
	restful.Add(service)
}
