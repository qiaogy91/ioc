package http

import "github.com/qiaogy91/ioc"

const (
	AppName = "http"
)

func Get() *Http {
	return ioc.Config().Get(AppName).(*Http)
}
