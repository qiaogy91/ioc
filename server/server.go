package server

import (
	"context"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/grpc"
	"github.com/qiaogy91/ioc/config/http"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
	"os"
	"os/signal"
	"syscall"
)

const (
	AppName = "application-server"
)

var configReq = &ioc.LoadConfigReq{
	ConfigFile: &ioc.ConfigFile{
		Enabled: true,
		Path:    "etc/application.yaml",
	},
	ConfigEnv: &ioc.ConfigEnv{
		Enabled: false,
		Prefix:  "",
	},
}

type Server struct {
	http   *http.Http
	grpc   *grpc.Server
	ch     chan os.Signal
	log    *zerolog.Logger
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

	s.log.Info().Msgf("loaded configs: %s", ioc.Config().List())
	s.log.Info().Msgf("loaded defaults: %s", ioc.Default().List())
	s.log.Info().Msgf("loaded controllers: %s", ioc.Controller().List())
	s.log.Info().Msgf("loaded apis: %s", ioc.Api().List())

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
			s.log.Info().Msgf("receive signal '%v', start graceful shutdown", v.String())

			if s.grpc.Enable {
				if err := s.grpc.Stop(s.ctx); err != nil {
					s.log.Error().Msgf("grpc graceful shutdown err: %s, force exit", err)
				}
			}

			if s.http.Enable {
				if err := s.http.Stop(s.ctx); err != nil {
					s.log.Error().Msgf("http graceful shutdown err: %s, force exit", err)
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
