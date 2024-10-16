package server

import (
	"context"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/grpc"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	AppName     = "Ioc"
	StopTimeout = 10 * time.Second
)

var configReq = "etc/application.yaml"

type Server struct {
	http *http.Http
	grpc *grpc.Server
	ch   chan os.Signal
	log  *slog.Logger
}

func (s *Server) Run(ctx context.Context) error {
	// 处理信号量
	s.ch = make(chan os.Signal, 1)
	signal.Notify(s.ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	s.http = http.Get()
	s.grpc = grpc.Get()

	s.log = log.Sub(AppName)

	s.log.Info("init success", slog.Any("config", ioc.Config().List()))
	s.log.Info("init success", slog.Any("default", ioc.Default().List()))
	s.log.Info("init success", slog.Any("controller", ioc.Controller().List()))
	s.log.Info("init success", slog.Any("api", ioc.Api().List()))

	go s.http.Start(ctx)
	go s.grpc.Start(ctx)

	s.waitSign()
	return nil
}

func (s *Server) waitSign() {
	for sg := range s.ch {
		switch v := sg.(type) {
		default:
			ctx, cancel := context.WithTimeout(context.Background(), StopTimeout)

			// 遍历每个名称空间，执行所有对象的Close 方法
			s.log.Warn("receive a signal", slog.String("signal", v.String()))
			if err := ioc.GetContainer().Close(ctx); err != nil {
				s.log.Error("close error", slog.String("reason", err.Error()))
			}

			// 清理资源
			cancel()
			s.log.Info("shutdown complete")
			return
		}
	}
}

func RunServ(ctx context.Context) error {
	ins := &Server{}
	return ins.Run(ctx)
}
