package grpc

import "github.com/qiaogy91/ioc"

const (
	AppName = "grpc"
)

func Get() *Server { return ioc.Config().Get(AppName).(*Server) }
