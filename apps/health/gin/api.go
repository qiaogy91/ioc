package gin

import (
	"github.com/qiaogy91/ioc"
	iocgin "github.com/qiaogy91/ioc/config/gin"
	iochttp "github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
)

const AppName = "health"

type Handler struct {
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 402 }

func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 路由注册
	r := iocgin.ModuleRouter(h)
	r.GET("", h.HealthHandler)

	h.log.Info().Msgf("Get the Health using http://%s/%s", iochttp.Get().Addr(), h.Name())
}

func init() {
	ioc.Api().Registry(&Handler{})
}
