package otlp

import (
	"context"
	"errors"
	"github.com/qiaogy91/ioc"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
)

type Impl struct {
	ioc.ObjectImpl
	log         *zerolog.Logger
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Endpoint    string `json:"endpoint" yaml:"endpoint"`
	Insecure    bool   `json:"insecure" yaml:"insecure"`
	shutdownFns []func(ctx context.Context) error
}

func (i *Impl) Name() string  { return AppName }
func (i *Impl) Priority() int { return 103 }
func (i *Impl) Close(ctx context.Context) error {
	var err error
	for _, fn := range i.shutdownFns {
		err = errors.Join(err, fn(ctx))
	}
	return err
}

func (i *Impl) Init() {
	// 注册全局 provider
	otel.SetTracerProvider(i.newTraceProvider())
	otel.SetMeterProvider(i.newMeterProvider())
	global.SetLoggerProvider(i.newLoggerProvider())

	// 定义上下文传播方式：传递额外的上下文信息、准的 context 上下文信息
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.Baggage{},
			propagation.TraceContext{},
		),
	)

}

func init() {
	ioc.Config().Registry(&Impl{shutdownFns: make([]func(ctx context.Context) error, 0)})
}
