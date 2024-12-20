package gin

import (
	"context"
	"fmt"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gin"
	iochttp "github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"log/slog"
)

const AppName = "metrics"

var (
	ins = &Handler{
		RequestTotalName:     "http_request_total",
		RequestHistogramName: "http_request_duration_histogram",
	}
)

type Handler struct {
	ioc.ObjectImpl
	log                          *slog.Logger
	RequestTotalName             string                  `json:"requestTotalName" yaml:"requestTotalName"`         // Counter 指标名称
	RequestHistogramName         string                  `json:"requestHistogramName" yaml:"requestHistogramName"` // Histogram 指标名称
	HttpRequestTotal             metric.Int64Counter     // 指标值
	HttpRequestDurationHistogram metric.Float64Histogram // 指标值
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 401 }
func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 从全局Provider 中获取一个 meter impl 来注册metric 指标
	ioc.OtlpMustEnabled()

	meter := otel.Meter(AppName)
	h.MetricRegistry(meter)

	gin.RootRouter().Use(h.MetricMiddleware) // 加载指标中间件，用来更新metric 指标值

	// 将 OpenTelemetry Metric 指标暴露到restful api 中
	r := gin.ModuleRouter(h)
	r.GET("", h.MetricHandler)

	h.log.Debug("Metric enabled",
		slog.String("visit", fmt.Sprintf("http://%s/%s", iochttp.Get().PrettyAddr(), AppName)))
}

func (h *Handler) Close(ctx context.Context) error {
	h.log.Debug("closed completed", slog.String("namespace", ioc.ApiNamespace))
	return nil
}

func init() {
	ioc.Api().Registry(ins)
}
