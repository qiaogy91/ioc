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
	tags := []string{"API 文档"}
	ws := gorestful.ModuleWebservice(h)
	ws.Route(ws.GET("doc.json").To(h.dockJson).
		Doc("查询文档信息").
		Metadata(restfulspec.KeyOpenAPITags, tags))

	ws.Route(ws.GET("doc.ui").To(h.docUI).
		Doc("查询文档面板").
		Metadata(restfulspec.KeyOpenAPITags, tags))

	h.log.Debug(fmt.Sprintf("Get the API doc using http://%s/%s/%s", http.Get().PrettyAddr(), AppName, "doc.json "))
	h.log.Debug(fmt.Sprintf("Get the Redoc UI using http://%s/%s/%s", http.Get().PrettyAddr(), AppName, "doc.ui "))
}

func (h *Handler) Close(ctx context.Context) error {
	h.log.Debug("closed completed", slog.String("namespace", ioc.ApiNamespace))
	return nil
}

func init() {
	ioc.Api().Registry(&Handler{})
}
