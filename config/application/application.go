package application

import (
	"github.com/qiaogy91/ioc"
)

var _ ioc.ObjectInterface = &Application{}

type Application struct {
	ioc.ObjectImpl
	AppName        string `json:"name" yaml:"name" toml:"name" env:"NAME"`
	AppDescription string `json:"description" yaml:"description" toml:"description" env:"DESCRIPTION"`
	Domain         string `json:"domain" yaml:"domain" toml:"domain" env:"DOMAIN"`
}

func (a *Application) Name() string {
	return AppName
}

func (a *Application) Priority() int {
	return 101
}

func init() {
	ioc.Config().Registry(&Application{})
}
