package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
	"net"
	"net/http"
	"time"
)

var (
	_   ioc.ObjectInterface = &Http{}
	ins                     = &Http{
		Host:              "127.0.0.1",
		Port:              8080,
		GinMode:           "debug",
		ReadHeaderTimeout: 30,
		ReadTimeout:       60,
		WriteTimeout:      60,
		IdleTimeout:       300,
		MaxHeaderSize:     "16kb",
		Otlp:              false,
	}
)

type Http struct {
	Host              string `json:"host" yaml:"host" env:"HOST"`
	Port              int    `json:"port" yaml:"port" env:"PORT"`
	GinMode           string `json:"ginMode" yaml:"ginMode"`                                               // 针对Gin 框架的模式
	ReadHeaderTimeout int    `json:"readHeaderTimeout" yaml:"readHeaderTimeout" env:"READ_HEADER_TIMEOUT"` // 读取请求头超时时间
	ReadTimeout       int    `json:"readTimeout" yaml:"readTimeout" env:"READ_TIMEOUT"`                    // 读取整个HTTP 的超时时间
	WriteTimeout      int    `json:"writeTimeout" yaml:"writeTimeout" env:"WRITE_TIMEOUT"`                 // 响应的超时时间
	IdleTimeout       int    `json:"idleTimeout" yaml:"idleTimeout" env:"IDLE_TIMEOUT"`                    // 开启Keepalive后，复用TCP 链接的超时时间
	MaxHeaderSize     string `json:"maxHeaderSize" yaml:"maxHeaderSize" env:"MAX_HEADER_SIZE"`             // HEADER 最大大小
	Otlp              bool   `json:"otlp" yaml:"otlp"`                                                     // 是否开启Trace

	ioc.ObjectImpl
	log            *slog.Logger
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
	h.log.Debug("HttpServer started",
		slog.String("listen", h.Addr()),
		slog.String("visit", fmt.Sprintf("http://%s", h.PrettyAddr())))

	if err := h.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return
		}
		h.log.Error("HttpServer Listen err", slog.Any("err", err))
	}
}
func (h *Http) Close(ctx context.Context) error {
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http graceful shutdown timeout, force exit")
	}
	h.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func (h *Http) PrettyAddr() string {
	if h.Host == "0.0.0.0" {
		// 如果用户配置的是 0.0.0.0 则从本地接口随便取出一个地址
		inters, err := net.InterfaceAddrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range inters {
			// 获取 IP 地址
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() {
				continue
			}
			// 打印非回环的 IPv4 地址
			if ip := ipNet.IP.To4(); ip != nil {
				return fmt.Sprintf("%s:%d", ip, h.Port)
			}
		}
	}
	return h.Addr()
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (h *Http) SetRouter(r http.Handler) {
	h.router = r
}

func init() {
	ioc.Config().Registry(ins)
}
