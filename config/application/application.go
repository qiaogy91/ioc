package application

import (
	"github.com/qiaogy91/ioc"
)

var _ ioc.ObjectInterface = &Application{}

type Application struct {
	ioc.ObjectImpl
	AppName        string `json:"appName" yaml:"appName"`
	AppDescription string `json:"description" yaml:"description"`
	Domain         string `json:"domain" yaml:"domain"`
}

func (a *Application) Name() string { return AppName }
func (a *Application) Priority() int {
	return 101
}
func (a *Application) ApplicationName() string { return a.AppName }
func init() {
	ioc.Config().Registry(&Application{})
}
