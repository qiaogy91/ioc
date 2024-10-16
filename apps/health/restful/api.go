package restful

import (
	"context"
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	iochttp "github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

const AppName = "health"

type Handler struct {
	ioc.ObjectImpl
	log  *slog.Logger
	Path string `json:"path" yaml:"path"`
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 402 }

func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 路由注册
	tags := []string{"健康检查"}
	ws := gorestful.ModuleWebservice(h)
	ws.Route(ws.GET("").To(h.HealthHandler).
		Doc("查询服务当前状态").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)

	h.log.Debug("Health enabled",
		slog.String("visit", fmt.Sprintf("http://%s/%s", iochttp.Get().PrettyAddr(), AppName)))
}

func (h *Handler) Close(ctx context.Context) error {
	h.log.Debug("closed completed", slog.String("namespace", ioc.ApiNamespace))
	return nil
}
func init() {
	ioc.Api().Registry(&Handler{})
}
