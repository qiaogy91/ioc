package trace

import (
	"context"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/application"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Trace struct {
	ioc.ObjectImpl
	log      *zerolog.Logger
	Enable   bool   `json:"enable" yaml:"enable"`
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Insecure bool   `json:"insecure" yaml:"insecure"`
	provider *trace.TracerProvider
}

func (t *Trace) Name() string  { return AppName }
func (t *Trace) Priority() int { return 103 }

func (t *Trace) options() (opts []otlptracehttp.Option) {
	if t.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	opts = append(opts, otlptracehttp.WithEndpoint(t.Endpoint))
	return
}

func (t *Trace) Init() {
	t.log = log.Sub(AppName)
	if !t.Enable {
		return
	}

	// resource 标注当前服务
	res, _ := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(application.Get().ApplicationName())),
	)
	// exporter
	exporter, err := otlptracehttp.New(context.Background(), t.options()...)
	if err != nil {
		panic(err)
	}
	// provider
	provider := trace.NewTracerProvider(
		trace.WithResource(res),                                        // resource
		trace.WithSampler(trace.AlwaysSample()),                        // sample
		trace.WithSpanProcessor(trace.NewBatchSpanProcessor(exporter)), // span processor
	)

	otel.SetTracerProvider(provider)                                                                                        // 注册到全局 global
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})) // 定义上下文传播方式：传递额外的上下文信息、准的 context 上下文信息
	t.provider = provider
}

func (t *Trace) Close(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

func init() {
	ioc.Config().Registry(&Trace{})
}
