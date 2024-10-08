package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	iochttp "github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
)

const AppName = "health"

type Handler struct {
	ioc.ObjectImpl
	log  *zerolog.Logger
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

	h.log.Info().Msgf("Get the Health using http://%s/%s", iochttp.Get().Addr(), h.Name())
}

func init() {
	ioc.Api().Registry(&Handler{})
}
