package trace

import "github.com/qiaogy91/ioc"

const (
	AppName = "trace"
)

func Get() *Trace {
	return ioc.Config().Get(AppName).(*Trace)
}
