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

var (
	_   ioc.ObjectInterface = &Server{}
	ins                     = &Server{
		Enabled: false,
		Host:    "127.0.0.1",
		Port:    18080,
		Token:   "my_secret_token",
	}
)

type Server struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Host    string `json:"host" yaml:"host"`
	Port    int    `json:"port" yaml:"port"`
	Token   string `json:"token" yaml:"token"`
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

func (s *Server) PrettyAddr() string {
	if s.Host == "0.0.0.0" {
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
				return fmt.Sprintf("%s:%d", ip, s.Port)
			}
		}
	}
	return s.Addr()
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

	s.log.Debug("GrpcServer started",
		slog.String("listen", s.Addr()),
		slog.String("visit", fmt.Sprintf("grpc://%s", s.PrettyAddr())))
	if err := s.server.Serve(lis); err != nil {
		s.log.Error("GrpcServer serve err", slog.Any("err", err))
	}
}

func (s *Server) Close(ctx context.Context) error {
	s.server.GracefulStop()
	s.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func (s *Server) Server() *grpc.Server {
	if s.server == nil {
		panic("GrpcServer server not initital")
	}
	return s.server
}

func init() {
	ioc.Config().Registry(ins)
}
