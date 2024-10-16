package grpc

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
)

const (
	AppName = config.GrpcName
)

func Get() *Server { return ioc.Config().Get(AppName).(*Server) }
