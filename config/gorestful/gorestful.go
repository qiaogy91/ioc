package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

type Framework struct {
	ioc.ObjectImpl
	Container *restful.Container
	log       *zerolog.Logger
}

func (f *Framework) Priority() int {
	return 104
}

func (f *Framework) Name() string {
	return AppName
}

func (f *Framework) Init() {
	f.log = log.Sub(AppName)
	f.Container = restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// 注册给Http服务器
	serv := http.Get()
	serv.SetRouter(f.Container)

	// 开启Trace
	if serv.Trace && trace.Get().Enable {
		f.Container.Filter(otelrestful.OTelFilter(application.Get().ApplicationName()))
		f.log.Info().Msg("restful trace enabled")
	}
}

func init() {
	ioc.Config().Registry(&Framework{})
}
