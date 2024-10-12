package log

import (
	"github.com/qiaogy91/ioc"
	"log/slog"
)

const (
	AppName      = "log"
	SubLoggerKey = "component"
)

func Sub(name string) *slog.Logger {
	return ioc.Config().Get(AppName).(*Logger).SubLogger(name)
}

func TextHandler() *slog.TextHandler {
	return ioc.Config().Get(AppName).(*Logger).HandlerConsole()
}
