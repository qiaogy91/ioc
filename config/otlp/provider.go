package otlp

import (
	"context"
	"github.com/qiaogy91/ioc/config/application"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/log/global"
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

func (i *Impl) getHttpExporter(ctx context.Context) (t trace.SpanExporter, m metric.Exporter, l log.Exporter, p *prometheus.Exporter) {
	var (
		traceOpts  []otlptracehttp.Option
		metricOpts []otlpmetrichttp.Option
		logOpts    []otlploghttp.Option
		err        error
	)
	// options
	traceOpts = append(traceOpts, otlptracehttp.WithEndpoint(i.HttpEndpoint))
	metricOpts = append(metricOpts, otlpmetrichttp.WithEndpoint(i.HttpEndpoint))
	logOpts = append(logOpts, otlploghttp.WithEndpoint(i.HttpEndpoint))
	if i.Insecure {
		traceOpts = append(traceOpts, otlptracehttp.WithInsecure())
		metricOpts = append(metricOpts, otlpmetrichttp.WithInsecure())
		logOpts = append(logOpts, otlploghttp.WithInsecure())
	}
	// exporters
	t, err = otlptracehttp.New(ctx, traceOpts...)
	m, err = otlpmetrichttp.New(ctx, metricOpts...)
	l, err = otlploghttp.New(ctx, logOpts...)
	p, err = prometheus.New()

	if err != nil {
		panic(err)
	}
	return
}

func (i *Impl) getGrpcExporter(ctx context.Context) (t trace.SpanExporter, m metric.Exporter, l log.Exporter, p *prometheus.Exporter) {
	var (
		traceOpts  []otlptracegrpc.Option
		metricOpts []otlpmetricgrpc.Option
		logOpts    []otlploggrpc.Option
		err        error
	)
	// options
	traceOpts = append(traceOpts, otlptracegrpc.WithEndpoint(i.GrpcEndpoint))
	metricOpts = append(metricOpts, otlpmetricgrpc.WithEndpoint(i.GrpcEndpoint))
	logOpts = append(logOpts, otlploggrpc.WithEndpoint(i.GrpcEndpoint))
	if i.Insecure {
		traceOpts = append(traceOpts, otlptracegrpc.WithInsecure())
		metricOpts = append(metricOpts, otlpmetricgrpc.WithInsecure())
		logOpts = append(logOpts, otlploggrpc.WithInsecure())
	}
	// exporters
	t, err = otlptracegrpc.New(ctx, traceOpts...)
	m, err = otlpmetricgrpc.New(ctx, metricOpts...)
	l, err = otlploggrpc.New(ctx, logOpts...)
	p, err = prometheus.New()

	if err != nil {
		panic(err)
	}
	return
}

func (i *Impl) RegistryProvider(ctx context.Context) {
	var (
		t trace.SpanExporter
		m metric.Exporter
		l log.Exporter
		p *prometheus.Exporter
	)

	switch i.GrpcEndpoint != "" {
	case true:
		t, m, l, p = i.getGrpcExporter(ctx)

	case false:
		t, m, l, p = i.getHttpExporter(ctx)
	}

	// trace
	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.TraceIDRatioBased(0.5)), // 采样率
		trace.WithBatcher(t),                            // 批量导出
		trace.WithResource(i.newResource()),             // 资源信息
	)
	// metric
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(m, metric.WithInterval(30*time.Second))), // 周期性导出
		metric.WithReader(p), // prometheus 导出
		metric.WithResource(i.newResource()),
	)
	// log
	logProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(l)), // 批量导出
		log.WithResource(i.newResource()),           // 资源信息
	)
	i.shutdownFns = append(i.shutdownFns, traceProvider.Shutdown, logProvider.Shutdown, meterProvider.Shutdown)

	otel.SetTracerProvider(traceProvider)
	global.SetLoggerProvider(logProvider)
	otel.SetMeterProvider(meterProvider)
}
