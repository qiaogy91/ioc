package otlp

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
)

const AppName = config.OtlpName

func Get() *Impl { return ioc.Config().Get(AppName).(*Impl) }
