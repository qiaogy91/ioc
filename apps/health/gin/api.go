package gin

import (
	"context"
	"fmt"
	"github.com/qiaogy91/ioc"
	iocgin "github.com/qiaogy91/ioc/config/gin"
	iochttp "github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

const AppName = "health"

type Handler struct {
	ioc.ObjectImpl
	log *slog.Logger
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 402 }

func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 路由注册
	r := iocgin.ModuleRouter(h)
	r.GET("", h.HealthHandler)

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
