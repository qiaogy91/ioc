package otlp

import (
	"context"
	"errors"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"log/slog"
)

var (
	ins = &Impl{
		HttpEndpoint: "127.0.0.1:4318",
		GrpcEndpoint: "127.0.0.1:4317",
		EnableTLS:    false,
		shutdownFns:  make([]func(ctx context.Context) error, 0),
	}
)

type Impl struct {
	ioc.ObjectImpl
	log          *slog.Logger
	HttpEndpoint string `json:"httpEndpoint" yaml:"httpEndpoint"`
	GrpcEndpoint string `json:"grpcEndpoint" yaml:"grpcEndpoint"`
	EnableTLS    bool   `json:"enableTLS" yaml:"enableTLS"`
	shutdownFns  []func(ctx context.Context) error
}

func (i *Impl) Name() string  { return AppName }
func (i *Impl) Priority() int { return 103 }
func (i *Impl) Close(ctx context.Context) error {
	var err error
	for _, fn := range i.shutdownFns {
		err = errors.Join(err, fn(ctx))
	}
	i.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return err
}

func (i *Impl) Init() {
	i.log = log.Sub(AppName)
	// 注册全局 impl
	i.RegistryProvider(context.Background())

	// 定义上下文传播方式：传递额外的上下文信息、准的 context 上下文信息
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)
	i.log.Debug("OpenTelemetry impl registered")
}

func init() {
	ioc.Config().Registry(ins)
}
