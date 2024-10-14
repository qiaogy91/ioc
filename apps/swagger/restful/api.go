package restful

import (
	"context"
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

const AppName = "swagger"

type Handler struct {
	ioc.ObjectImpl
	log *slog.Logger
}

func (h *Handler) Name() string { return AppName }
func (h *Handler) Priority() int {
	return 403
}

func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 路由注册
	ws := gorestful.ModuleWebservice(h)
	ws.Route(
		ws.GET("doc.json").To(h.restfulSwagger).
			Doc("查询文档信息").
			Metadata(restfulspec.KeyOpenAPITags, []string{"API 文档"}),
	)

	h.log.Info(fmt.Sprintf("Get the API doc using http://%s/%s/%s", http.Get().Addr(), AppName, "doc.json "))
}
func (h *Handler) Close(ctx context.Context) error {
	h.log.Info("closed completed", slog.String("namespace", ioc.ApiNamespace))
	return nil
}

func init() {
	ioc.Api().Registry(&Handler{})
}
