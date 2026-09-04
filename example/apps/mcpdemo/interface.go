package mcpdemo

import (
	"context"
	"github.com/qiaogy91/ioc"
)

const (
	AppName        = "mcpdemo"
	HandlerName    = "handler_mcpdemo"
	ControllerName = "controller_mcpdemo"
)

func GetController() Controller {
	return ioc.Controller().Get(ControllerName).(Controller)
}

type Controller interface {
	GetUser(ctx context.Context, id int) (*User, error) // GetUser 根据 ID 获取用户信息
}
