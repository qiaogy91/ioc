package log

import (
	"github.com/qiaogy91/ioc"
	"github.com/rs/zerolog"
)

const (
	AppName      = "log_bak"
	SubLoggerKey = "component"
)

func Sub(name string) *zerolog.Logger {
	return ioc.Config().Get(AppName).(*Logger).Logger(name)
}
