package http

import (
	"context"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

var _ ioc.ObjectInterface = &Http{}

type Http struct {
	Enable            bool   `json:"enable" yaml:"enable" env:"ENABLE"`
	Host              string `json:"host" yaml:"host" env:"HOST"`
	Port              int    `json:"port" yaml:"port" env:"PORT"`
	ReadHeaderTimeout int    `json:"readHeaderTimeout" yaml:"readHeaderTimeout" env:"READ_HEADER_TIMEOUT"` // 读取请求头超时时间
	ReadTimeout       int    `json:"readTimeout" yaml:"readTimeout" env:"READ_TIMEOUT"`                    // 读取整个HTTP 的超时时间
	WriteTimeout      int    `json:"writeTimeout" yaml:"writeTimeout" env:"WRITE_TIMEOUT"`                 // 响应的超时时间
	IdleTimeout       int    `json:"idleTimeout" yaml:"idleTimeout" env:"IDLE_TIMEOUT"`                    // 开启Keepalive后，复用TCP 链接的超时时间
	MaxHeaderSize     string `json:"maxHeaderSize" yaml:"maxHeaderSize" env:"MAX_HEADER_SIZE"`             // HEADER 最大大小
	Trace             bool   `json:"trace" yaml:"trace"`                                                   // 是否开启Trace

	ioc.ObjectImpl
	log            *zerolog.Logger
	router         http.Handler
	server         *http.Server
	maxHeaderBytes uint64 // 解析后的数据
}

func (h *Http) Name() string {
	return AppName
}

func (h *Http) Priority() int {
	return 108
}

func (h *Http) Init() {
	mhz, err := humanize.ParseBytes(h.MaxHeaderSize)
	if err != nil {
		panic(err)
	}
	h.maxHeaderBytes = mhz   // 最大请求头
	h.log = log.Sub(AppName) // 日志句柄

	h.server = &http.Server{
		ReadHeaderTimeout: time.Duration(h.ReadHeaderTimeout) * time.Second,
		ReadTimeout:       time.Duration(h.ReadTimeout) * time.Second,
		WriteTimeout:      time.Duration(h.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(h.IdleTimeout) * time.Second,
		MaxHeaderBytes:    int(h.maxHeaderBytes),
		Addr:              h.Addr(),
		Handler:           h.router,
	}
}

func (h *Http) Start(ctx context.Context) {
	h.log.Info().Msgf("Started HttpServer at: %s", h.Addr())
	if err := h.server.ListenAndServe(); err != nil {
		h.log.Error().Msg(err.Error())
	}
}
func (h *Http) Stop(ctx context.Context) error {
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http graceful shutdown timeout, force exit")
	}
	h.log.Info().Msg("Shutdown HttpServer Complete")
	return nil
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (h *Http) SetRouter(r http.Handler) {
	h.router = r
}

func init() {
	ioc.Config().Registry(&Http{})
}
