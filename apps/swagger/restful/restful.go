package restful

import (
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/qiaogy91/ioc/apps/swagger"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/http"
	"log/slog"
)

func (h *Handler) docUI(r *restful.Request, w *restful.Response) {
	w.Header().Set("Content-Type", "text/html")

	docApi := fmt.Sprintf("http://%s/%s/%s", http.Get().PrettyAddr(), AppName, "doc.json ")
	docHtml := fmt.Sprintf(swagger.DocHtml, docApi)

	if _, err := w.Write([]byte(docHtml)); err != nil {
		h.log.Error("swagger writeAsJson failed", slog.Any("err", err))
	}
}

// BuildSwagger 定义swagger 配置
func (h *Handler) BuildSwagger() restfulspec.Config {
	return restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		PostBuildSwaggerObjectHandler: func(s *spec.Swagger) {
			s.Info = &spec.Info{
				InfoProps: spec.InfoProps{
					Title:       application.Get().ApplicationName(),
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

func (h *Handler) dockJson(request *restful.Request, response *restful.Response) {
	swg := restfulspec.BuildSwagger(h.BuildSwagger())
	if err := response.WriteAsJson(swg); err != nil {
		h.log.Error("swagger writeAsJson failed", slog.Any("err", err))
	}
}
