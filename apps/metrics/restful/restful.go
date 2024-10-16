package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qiaogy91/ioc/config/application"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"log/slog"
	"time"
)

func (h *Handler) MetricHandler(req *restful.Request, rsp *restful.Response) {
	promhttp.Handler().ServeHTTP(rsp, req.Request)
}

// MetricRegistry 创建指标
func (h *Handler) MetricRegistry(meter metric.Meter) {
	var err error
	h.HttpRequestDurationHistogram, err = meter.Float64Histogram(
		h.RequestHistogramName,
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("{seconds}"),
	)

	h.HttpRequestTotal, err = meter.Int64Counter(
		h.RequestTotalName,
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{call}"),
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
		attribute.String("service", application.Get().AppName),
		attribute.String("method", req.Request.Method),
		attribute.String("path", req.SelectedRoutePath())),
	)
	h.log.Info("延迟记录", slog.Float64("duration", float64(time.Since(start).Milliseconds())))
	h.HttpRequestDurationHistogram.Record(req.Request.Context(), float64(time.Since(start).Milliseconds()), metric.WithAttributes(
		attribute.String("service", application.Get().AppName),
		attribute.String("method", req.Request.Method),
		attribute.String("path", req.SelectedRoutePath()),
	))
}
