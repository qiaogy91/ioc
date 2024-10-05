package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/qiaogy91/ioc/config/application"
)

// BuildSwagger 定义swagger 配置
func (h *Handler) BuildSwagger() restfulspec.Config {
	return restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       application.Get().AppName,
					Description: application.Get().AppDescription,
					License: &spec.License{
						LicenseProps: spec.LicenseProps{
							Name: "MIT",
							URL:  "https://opensource.org/licenses/MIT",
						},
					},
				},
			}
		},
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
	}
}

func (h *Handler) restfulSwagger(request *restful.Request, response *restful.Response) {
	swagger := restfulspec.BuildSwagger(h.BuildSwagger())
	if err := response.WriteAsJson(swagger); err != nil {
		h.log.Error().Msgf("swagger writeAsJson failed, %s", err)
	}
}
