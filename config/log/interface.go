package log

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
	"log/slog"
)

const (
	AppName      = config.LogName
	SubLoggerKey = "component"
)

func Sub(name string) *slog.Logger {
	return ioc.Config().Get(AppName).(*Logger).SubLogger(name)
}

func TextHandler() *slog.TextHandler {
	return ioc.Config().Get(AppName).(*Logger).HandlerConsole()
}
