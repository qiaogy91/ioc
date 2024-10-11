package otlp

import (
	"context"
	"github.com/qiaogy91/ioc/config/application"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"time"
)

// Resource
func (i *Impl) newResource() *resource.Resource {
	// 自动获取Resource 属性
	autoRes, err := resource.New(
		context.Background(),
		resource.WithHost(),         // 自动检测主机信息
		resource.WithProcess(),      // 自动检测进程信息
		resource.WithTelemetrySDK(), // 自动检测 SDK 信息
	)
	if err != nil {
		panic(err)
	}

	// 自定义Resource 属性
	customRes := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(application.Get().ApplicationName()), // 自定义服务名
	)

	// 最终使用合并后的Resource
	finalRes, err := resource.Merge(autoRes, customRes)
	if err != nil {
		panic(err)
	}
	return finalRes
}

// Exporter
func (i *Impl) newTracerExporter(ctx context.Context) *otlptrace.Exporter {
	// opts
	var opts []otlptracehttp.Option
	opts = append(opts, otlptracehttp.WithEndpoint(i.Endpoint))
	if i.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	// exporter
	exp, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return exp
}

func (i *Impl) newMeterExporter(ctx context.Context) (*otlpmetrichttp.Exporter, *prometheus.Exporter) {
	// opts
	var opts []otlpmetrichttp.Option
	if i.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}
	opts = append(opts, otlpmetrichttp.WithEndpoint(i.Endpoint))

	// otlp exporter
	otlpExp, err := otlpmetrichttp.New(ctx, opts...)
	if err != nil {
		panic(err)
	}
	// prom exporter
	promExp, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	return otlpExp, promExp
}

func (i *Impl) newLoggerExporter(ctx context.Context) *otlploghttp.Exporter {
	// opts
	var opts []otlploghttp.Option
	if i.Insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}
	opts = append(opts, otlploghttp.WithEndpoint(i.Endpoint))

	// exporter
	exp, err := otlploghttp.New(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return exp
}

// Provider
func (i *Impl) newTraceProvider() *trace.TracerProvider {
	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(0.5)),              // 采样率
		trace.WithBatcher(i.newTracerExporter(context.Background())), // 批量导出
		trace.WithResource(i.newResource()),                          // 资源信息
	)

	i.shutdownFns = append(i.shutdownFns, provider.Shutdown)
	return provider
}

func (i *Impl) newMeterProvider() *metric.MeterProvider {
	otlpExp, promExp := i.newMeterExporter(context.Background())
	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(otlpExp, metric.WithInterval(30*time.Second))),
		metric.WithReader(promExp),
	)
	i.shutdownFns = append(i.shutdownFns, provider.Shutdown)
	return provider
}

func (i *Impl) newLoggerProvider() *log.LoggerProvider {
	provider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(i.newLoggerExporter(context.Background()))),
	)

	i.shutdownFns = append(i.shutdownFns, provider.Shutdown)
	return provider
}
