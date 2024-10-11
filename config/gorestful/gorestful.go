package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/config/otlp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
	"time"
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

	f.Container.Filter(f.AccessLog)

	// 开启Trace
	//if serv.Trace && trace.Get().Enable {
	//	f.Container.Filter(otelrestful.OTelFilter(application.Get().ApplicationName()))
	//	f.log.Info().Msg("restful trace enabled")
	//}
	// 替换为otlp trace
	if serv.Trace && otlp.Get().Enabled {
		f.Container.Filter(otelrestful.OTelFilter(application.Get().ApplicationName()))
		f.log.Info().Msg("restful trace enabled")
	}
}

func (f *Framework) AccessLog(r *restful.Request, w *restful.Response, chain *restful.FilterChain) {
	f.log = log.Sub("accessLog")
	start := time.Now()
	chain.ProcessFilter(r, w)

	// 返回Response 时记录日志
	f.log.Info().Msgf("%-20s | %-15s | %-5d | %-10s | %s",
		time.Since(start),
		r.Request.RemoteAddr,
		w.StatusCode(),
		r.Request.Method,
		r.Request.URL.Path)
}

func init() {
	ioc.Config().Registry(&Framework{})
}
