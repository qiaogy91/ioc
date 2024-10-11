package application

import "github.com/qiaogy91/ioc"

const (
	AppName = "application"
)

func Get() *Application {
	return ioc.Config().Get(AppName).(*Application)
}
