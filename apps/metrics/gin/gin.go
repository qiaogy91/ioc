package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/qiaogy91/ioc/config/application"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"time"
)

func (h *Handler) MetricHandler(ctx *gin.Context) {
	promhttp.Handler().ServeHTTP(ctx.Writer, ctx.Request)
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
func (h *Handler) MetricMiddleware(ctx *gin.Context) {
	start := time.Now()
	ctx.Next()

	h.HttpRequestTotal.Add(ctx.Request.Context(), 1, metric.WithAttributes(
		attribute.String("service", application.Get().AppName),
		attribute.String("method", ctx.Request.Method),
		attribute.String("path", ctx.FullPath())),
	)

	h.HttpRequestDurationHistogram.Record(ctx.Request.Context(), time.Since(start).Seconds(),
		metric.WithAttributes(
			attribute.String("service", application.Get().AppName),
			attribute.String("method", ctx.Request.Method),
			attribute.String("path", ctx.FullPath()),
		))
}
