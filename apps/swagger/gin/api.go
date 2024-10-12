package gin

import (
	"fmt"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gin"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	r := gin.ModuleRouter(h)
	r.GET("/doc.json", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h.log.Info(fmt.Sprintf("Get the API doc using http://%s/%s/%s", http.Get().Addr(), AppName, "doc.json"))
}

func init() {
	ioc.Api().Registry(&Handler{})
}
