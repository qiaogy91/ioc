package http

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
)

const (
	AppName = config.HttpName
)

func Get() *Http {
	return ioc.Config().Get(AppName).(*Http)
}
