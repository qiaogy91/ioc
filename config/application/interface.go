package application

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
)

const (
	AppName = config.ApplicationName
)

func Get() *Application {
	return ioc.Config().Get(AppName).(*Application)
}
