package gin

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
)

const AppName = "metrics"

type Handler struct {
	ioc.ObjectImpl
	log                     *zerolog.Logger
	ApiStatusEnable         bool      `json:"apiStatusEnable" yaml:"apiStatusEnable"`                 // 是否开启内置的 api status 指标统计
	RequestHistogram        bool      `json:"requestHistogram" yaml:"requestHistogram"`               // Histogram 类型标签，是否开启
	RequestHistogramName    string    `json:"requestHistogramName" yaml:"requestHistogramName"`       // Histogram 标签名称
	RequestHistogramBucket  []float64 `json:"requestHistogramBucket" yaml:"requestHistogramBucket"`   // Histogram bucket 边界
	RequestSummary          bool      `json:"requestSummary" yaml:"requestSummary"`                   // Summary 类型标签，是否开启
	RequestSummaryName      string    `json:"requestSummaryName" yaml:"requestSummaryName"`           // Summary 标签名称
	RequestSummaryObjective []float64 `json:"requestSummaryObjective" yaml:"requestSummaryObjective"` // Summary bucket 边界
	RequestTotal            bool      `json:"requestTotal" yaml:"requestTotal"`                       // Counter 类型标签，是否开启
	RequestTotalName        string    `json:"requestTotalName" yaml:"requestTotalName"`               // Counter 标签名称

	// 三个指标
	HttpRequestTotal             *prometheus.CounterVec   // 请求总数
	HttpRequestDurationHistogram *prometheus.HistogramVec // 请求时长柱状图
	HttpRequestDurationSummary   *prometheus.SummaryVec   // 请求时长分位数
}

func (h *Handler) Name() string  { return AppName }
func (h *Handler) Priority() int { return 401 }
func (h *Handler) Init() {
	h.log = log.Sub(AppName)

	// 如果未注册，则先创建 collector 采集器，并将采集器注册进注册表
	h.CollectInit()

	// 注册路由
	ws := gorestful.ModuleWebservice(h)
	ws.Route(ws.GET("").To(h.MetricHandler).
		Doc("指标暴露").
		Metadata(restfulspec.KeyOpenAPITags, []string{"指标监控"}),
	)

	// 加载到全局中间件
	gorestful.RootContainer().Filter(h.MetricMiddleware)
	h.log.Info().Msgf("Get the Metric using http://%s/%s", http.Get().Addr(), h.Name())
}

func init() {
	ioc.Api().Registry(&Handler{})
}
