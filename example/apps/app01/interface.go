package app01

import (
	"context"
	"github.com/qiaogy91/ioc"
)

const (
	AppName = "app01"
)

func GetSvc() Service {
	return ioc.Controller().Get(AppName).(Service)
}

type Service interface {
	CreatTable(ctx context.Context) error
	ServiceServer
}
