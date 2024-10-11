package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

func (h *Handler) MetricHandler(req *restful.Request, rsp *restful.Response) {
	promhttp.Handler().ServeHTTP(rsp, req.Request)
}

// MetricRegistry 创建指标
func (h *Handler) MetricRegistry() {
	var err error
	var meter = otel.Meter(AppName)
	h.HttpRequestDurationHistogram, err = meter.Float64Histogram(
		h.RequestHistogramName,
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("s"),
	)

	h.HttpRequestTotal, err = meter.Int64Counter(
		h.RequestTotalName,
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("count"),
	)

	if err != nil {
		panic(err)
	}
}

// MetricMiddleware 指标中间件
func (h *Handler) MetricMiddleware(req *restful.Request, rsp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()

	chain.ProcessFilter(req, rsp)

	h.HttpRequestTotal.Add(req.Request.Context(), 1, metric.WithAttributes(
		attribute.String("method", req.Request.Method),
		attribute.String("path", req.SelectedRoutePath())),
	)

	h.HttpRequestDurationHistogram.Record(req.Request.Context(), time.Since(start).Seconds(), metric.WithAttributes(
		attribute.String("method", req.Request.Method),
		attribute.String("path", req.SelectedRoutePath()),
	))
}
