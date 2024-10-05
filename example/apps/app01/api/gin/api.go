package gin

import (
	"github.com/qiaogy91/ioc"
	iocgin "github.com/qiaogy91/ioc/config/gin"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/example/apps/app01"
	"github.com/rs/zerolog"
)

type Handler struct {
	ioc.ObjectImpl
	svc app01.Service
	log *zerolog.Logger
}

// Init
// 因为handler 在设置路由时必须获取根Root 后才能进行设置，而获取根Root 先要导入ioc 中对应gin、restful 的包
// 因此导入哪个包，就执行哪个框架的Init() 方法，这个方法会自动将handler 注册到 HTTP 框架上
// 由此来做出自动适配，即客户端写什么样的视图函数，自动启动什么样的框架
func (h *Handler) Init() {
	h.svc = app01.GetSvc()
	h.log = log.Sub(app01.AppName)

	// Gin 框架路由注册
	route := iocgin.ModuleRouter(h)
	route.GET("", h.ginList)
	route.POST("", h.ginCreate)
	route.POST("/table", h.ginCreatTable)
}
func (h *Handler) Name() string  { return app01.AppName }
func (h *Handler) Priority() int { return 401 }

func init() {
	ioc.Api().Registry(&Handler{})
}
