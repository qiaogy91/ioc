package application

import "github.com/qiaogy91/ioc"

const (
	AppName = "app"
)

func Get() *Application {
	return ioc.Config().Get(AppName).(*Application)
}
