package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

func (h *Handler) MetricHandler(ctx *gin.Context) {
	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
}

// Objectives 偏差计算
// 分位数0.5，计算 (1-0.5)*0.1 = 0.05
// 分位数0.9，计算 (1-0.9)*0.1 = 0.001
// 分位数0.99，计算 (1-0.99)*0.1 = 0.0001
func (h *Handler) Objectives() map[float64]float64 {
	objectives := map[float64]float64{}
	for _, v := range h.RequestSummaryObjective {
		objectives[v] = (1 - v) * 0.1
	}
	return objectives
}

// CollectInit 创建三种指标，并注册
func (h *Handler) CollectInit() {
	h.HttpRequestDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    h.RequestHistogramName,
			Help:    "Histogram of the duration of HTTP requests",
			Buckets: h.RequestHistogramBucket,
		},
		[]string{"method", "path"},
	)

	h.HttpRequestDurationSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       h.RequestSummaryName,
			Help:       "Histogram of the duration of HTTP requests",
			Objectives: h.Objectives(),
		},
		[]string{"method", "path"},
	)
	h.HttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: h.RequestTotalName,
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	// 注册进 registry
	prometheus.MustRegister(h) // 自身已经实现了Describe、Collect 方法，因此Handler 本身就是一个采集器
}

// Describe 获取指标描述
func (h *Handler) Describe(ch chan<- *prometheus.Desc) {
	if h.RequestHistogram {
		h.HttpRequestDurationHistogram.Describe(ch)
	}
	if h.RequestSummary {
		h.HttpRequestDurationSummary.Describe(ch)
	}
	if h.RequestTotal {
		h.HttpRequestTotal.Describe(ch)
	}
}

// Collect 获取指标值
func (h *Handler) Collect(ch chan<- prometheus.Metric) {
	if h.RequestHistogram {
		h.HttpRequestDurationHistogram.Collect(ch)
	}
	if h.RequestSummary {
		h.HttpRequestDurationSummary.Collect(ch)
	}
	if h.RequestTotal {
		h.HttpRequestTotal.Collect(ch)
	}
}

func (h *Handler) MetricMiddleware(ctx *gin.Context) {
	start := time.Now()

	ctx.Next()
	if h.RequestTotal {
		h.HttpRequestTotal.WithLabelValues(ctx.Request.Method, ctx.FullPath(), strconv.Itoa(ctx.Writer.Status())).Inc()
	}
	if h.RequestSummary {
		h.HttpRequestDurationSummary.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Observe(time.Since(start).Seconds())
	}
	if h.RequestHistogram {
		h.HttpRequestDurationHistogram.WithLabelValues(ctx.Request.Method, ctx.FullPath()).Observe(time.Since(start).Seconds())
	}
}
