package grpc

import (
	"context"
	"fmt"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log/slog"
	"net"
	"slices"
)

var _ ioc.ObjectInterface = &Server{}

type Server struct {
	Enable bool   `json:"enable" yaml:"enable"`
	Host   string `json:"host" yaml:"host"`
	Port   int    `json:"port" yaml:"port"`
	Token  string `json:"token" yaml:"token"`
	ioc.ObjectImpl
	server *grpc.Server
	log    *slog.Logger
}

func (s *Server) Name() string {
	return AppName
}

func (s *Server) Priority() int { return 107 }

func (s *Server) Init() {
	s.log = log.Sub(AppName)
	s.server = grpc.NewServer(grpc.UnaryInterceptor(s.TokenAuth))
}

func (s *Server) TokenAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 获取元数据
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Aborted, "未获取到meta 数据")
	}

	// token 校验
	token := md.Get("token")
	if !slices.Contains(token, s.Token) {
		return nil, status.Errorf(codes.PermissionDenied, "非法token")
	}

	return handler(ctx, req) // 通过则放行
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s *Server) Start(ctx context.Context) {
	lis, err := net.Listen("tcp", s.Addr())
	if err != nil {
		s.log.Error("GrpcServer listen err", slog.Any("err", err))
		return
	}

	s.log.Info("GrpcServer Started", slog.String("addr", s.Addr()))
	if err := s.server.Serve(lis); err != nil {
		s.log.Error("GrpcServer serve err", slog.Any("err", err))
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.server.GracefulStop()
	s.log.Info("GrpcServer Shutdown Complete")
	return nil
}

func (s *Server) Server() *grpc.Server {
	if s.server == nil {
		panic("GrpcServer server not initital")
	}
	return s.server
}

func init() {
	ioc.Config().Registry(&Server{})
}
