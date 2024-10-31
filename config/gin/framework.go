package gin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"log/slog"
)

var (
	_   ioc.ObjectInterface = &Framework{}
	ins                     = &Framework{
		Mode:   "release",
		Engine: nil,
		log:    nil,
	}
)

type Framework struct {
	ioc.ObjectImpl
	Mode   string `toml:"mode" json:"mode" yaml:"mode" env:"Mode"` // gin mode
	Engine *gin.Engine
	log    *slog.Logger
}

func (f *Framework) Name() string {
	return AppName
}

func (f *Framework) Priority() int {
	return 104
}

func (f *Framework) Close(ctx context.Context) error {
	f.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func (f *Framework) Init() {
	// 获取Http 服务器
	serv := http.Get()

	// 初始化当前实例
	gin.SetMode(serv.GinMode) // 设置Gin 全局模式（要在Engine 初始化之前完成，否则不生效）

	f.log = log.Sub(AppName)
	f.Engine = gin.Default()
	f.Engine.Use(gin.Recovery())

	// 开启Trace
	if serv.Otlp {
		ioc.OtlpMustEnabled()
		f.Engine.Use(otelgin.Middleware(application.Get().ApplicationName()))
		f.log.Debug("Gin trace enabled")
	}

	// 注册给Http服务器
	serv.SetRouter(f.Engine)

}

func init() {
	ioc.Config().Registry(ins)
}
