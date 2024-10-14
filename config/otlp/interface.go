package otlp

import "github.com/qiaogy91/ioc"

const AppName = "otlp"

func Get() *Impl { return ioc.Config().Get(AppName).(*Impl) }
