package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

var _ ioc.ObjectInterface = &Framework{}

type Framework struct {
	ioc.ObjectImpl
	Mode   string `toml:"mode" json:"mode" yaml:"mode" env:"Mode"` // gin mode
	Engine *gin.Engine
	log    *zerolog.Logger
}

func (f *Framework) Name() string {
	return AppName
}

func (f *Framework) Priority() int {
	return 104
}

func (f *Framework) Init() {
	f.log = log.Sub(AppName)
	f.Engine = gin.Default()
	f.Engine.Use(gin.Recovery())
	gin.SetMode(f.Mode)

	// 注册给Http服务器
	serv := http.Get()
	serv.SetRouter(f.Engine)

	// 开启Trace
	if serv.Trace && trace.Get().Enable {
		f.Engine.Use(otelgin.Middleware(application.Get().Name()))
		f.log.Info().Msg("gin trace enabled")
	}
}

func init() {
	ioc.Config().Registry(&Framework{})
}
