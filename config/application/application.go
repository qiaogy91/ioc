package application

import (
	"context"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

var (
	_   ioc.ObjectInterface = &Application{}
	ins                     = &Application{
		AppName:        "DefaultApp",
		AppDescription: "my default app service",
		Domain:         "example.com",
	}
)

type Application struct {
	ioc.ObjectImpl
	log            *slog.Logger
	AppName        string `json:"appName" yaml:"appName"`
	AppDescription string `json:"description" yaml:"description"`
	Domain         string `json:"domain" yaml:"domain"`
}

func (a *Application) Name() string { return AppName }
func (a *Application) Priority() int {
	return 102
}
func (a *Application) Init() {
	a.log = log.Sub(AppName)
}

func (a *Application) Close(ctx context.Context) error {
	a.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func (a *Application) ApplicationName() string { return a.AppName }
func init() {
	ioc.Config().Registry(ins)
}
