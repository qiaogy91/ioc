package restful

import (
	"context"
	"fmt"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/http"
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
	RequestTotalName             string                  `json:"requestTotalName" yaml:"requestTotalName"`         // Counter 标签名称
	RequestHistogramName         string                  `json:"requestHistogramName" yaml:"requestHistogramName"` // Histogram 标签名称
	HttpRequestTotal             metric.Int64Counter     // 请求总数
	HttpRequestDurationHistogram metric.Float64Histogram // 请求时长柱状图
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 401 }
func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	ioc.OtlpMustEnabled()

	// 从全局Provider 中获取一个 meter impl 来注册metric 指标
	meter := otel.Meter(AppName)
	h.MetricRegistry(meter)

	// 加载指标中间件，用来更新metric 指标值
	gorestful.RootContainer().Filter(h.MetricMiddleware)

	// 将 OpenTelemetry Metric 指标暴露到restful api 中
	ws := gorestful.ModuleWebservice(h)
	ws.Route(ws.GET("").To(h.MetricHandler).
		Doc("指标暴露").
		Metadata(restfulspec.KeyOpenAPITags, []string{"指标监控"}),
	)

	h.log.Debug("Metric enabled",
		slog.String("visit", fmt.Sprintf("http://%s/%s", http.Get().PrettyAddr(), AppName)))
}

func (h *Handler) Close(ctx context.Context) error {
	h.log.Debug("closed completed", slog.String("namespace", ioc.ApiNamespace))
	return nil
}

func init() {
	ioc.Api().Registry(ins)
}
