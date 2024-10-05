package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/example/apps/app01"
	"github.com/rs/zerolog"
)

type Handler struct {
	ioc.ObjectImpl
	svc app01.Service
	log *zerolog.Logger
}

func (h *Handler) Name() string  { return app01.AppName }
func (h *Handler) Priority() int { return 401 }

// Init
// 因为handler 在设置路由时必须获取根Root 后才能进行设置，而获取根Root 先要导入ioc 中对应gin、restful 的包
// 因此导入哪个包，就执行哪个框架的Init() 方法，这个方法会自动将handler 注册到 HTTP 框架上
// 由此来做出自动适配，即客户端写什么样的视图函数，自动启动什么样的框架
func (h *Handler) Init() {
	h.svc = app01.GetSvc()
	h.log = log.Sub(app01.AppName)

	// Restful 框架路由注册
	ws := gorestful.ModuleWebservice(h)
	ws.Route(ws.GET("").To(h.restfulList).Doc("用户列表"))
	ws.Route(ws.POST("").To(h.restfulCreate).Doc("创建用户"))
	ws.Route(ws.POST("/table").To(h.restfulCreateTable).Doc("创建表结构"))
	// 打印所有已注册的路由
	for _, ws := range restful.RegisteredWebServices() {
		for _, r := range ws.Routes() {
			h.log.Info().Msgf("%-6s %s", r.Method, r.Path)
		}
	}
}

func init() {
	ioc.Api().Registry(&Handler{})
}
