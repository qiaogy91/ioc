package restful

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
)

const AppName = "swagger"

type Handler struct {
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (h *Handler) Name() string { return AppName }
func (h *Handler) Priority() int {
	return 403
}

func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 路由注册
	ws := gorestful.ModuleWebservice(h)
	ws.Route(ws.GET("doc.json").To(h.restfulSwagger))

	h.log.Info().Msgf("Get the API doc using http://%s/%s/%s", http.Get().Addr(), AppName, "doc.json ")
}

func init() {
	ioc.Api().Registry(&Handler{})
}
