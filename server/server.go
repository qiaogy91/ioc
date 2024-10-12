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
)

const (
	AppName = "application-server"
)

var configReq = "etc/application.yaml"

type Server struct {
	http   *http.Http
	grpc   *grpc.Server
	ch     chan os.Signal
	log    *slog.Logger
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *Server) Run(ctx context.Context) error {
	// 初始化ioc
	err := ioc.ConfigIocObject(configReq)
	if err != nil {
		return err
	}

	// 处理信号量
	s.ch = make(chan os.Signal, 1)
	signal.Notify(s.ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	s.ctx, s.cancel = context.WithCancel(context.Background())

	s.http = http.Get()
	s.grpc = grpc.Get()

	s.log = log.Sub(AppName)

	s.log.Info("config namespace", slog.Any("loaded", ioc.Config().List()))
	s.log.Info("default namespace", slog.Any("loaded", ioc.Default().List()))
	s.log.Info("controller namespace", slog.Any("loaded", ioc.Controller().List()))
	s.log.Info("apis namespace", slog.Any("loaded", ioc.Api().List()))

	if s.http.Enable {
		go s.http.Start(ctx)
	}

	if s.grpc.Enable {
		go s.grpc.Start(ctx)
	}
	s.waitSign()
	return nil
}

func (s *Server) waitSign() {
	defer s.cancel()

	for sg := range s.ch {
		switch v := sg.(type) {
		default:
			s.log.Info("graceful shutdown", slog.String("reason", v.String()))

			if s.grpc.Enable {
				if err := s.grpc.Stop(s.ctx); err != nil {
					s.log.Error("graceful shutdown error", slog.Any("err", err))
				}
			}

			if s.http.Enable {
				if err := s.http.Stop(s.ctx); err != nil {
					s.log.Error("http graceful shutdown error", slog.Any("err", err))
				}
			}
			return
		}
	}
}

func RunServ(ctx context.Context) error {
	ins := &Server{}
	return ins.Run(ctx)
}
