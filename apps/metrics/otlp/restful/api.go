package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/metric"
)

const AppName = "metrics"

type Handler struct {
	ioc.ObjectImpl
	log                          *zerolog.Logger
	RequestTotalName             string                  `json:"requestTotalName" yaml:"requestTotalName"`             // Counter 标签名称
	RequestHistogramName         string                  `json:"requestHistogramName" yaml:"requestHistogramName"`     // Histogram 标签名称
	RequestHistogramBucket       []float64               `json:"requestHistogramBucket" yaml:"requestHistogramBucket"` // Histogram bucket 边界
	HttpRequestTotal             metric.Int64Counter     // 请求总数
	HttpRequestDurationHistogram metric.Float64Histogram // 请求时长柱状图
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 401 }
func (h *Handler) Init() {
	h.log = log.Sub(AppName)
	h.MetricRegistry() // 指标注册

	r := gorestful.ModuleWebservice(h)
	r.Route(r.GET("").To(h.MetricHandler).
		Doc("指标暴露").
		Metadata(restfulspec.KeyOpenAPITags, []string{"指标监控"}),
	)

	gorestful.RootContainer().Filter(h.MetricMiddleware)
	h.log.Info().Msgf("Get the Metric using http://%s/%s", http.Get().Addr(), h.Name())
}

func init() {
	ioc.Api().Registry(&Handler{})
}
