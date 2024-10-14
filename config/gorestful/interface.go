package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc"
	"path"
)

const (
	AppName = "restful"
)

func RootContainer() *restful.Container {
	return ioc.Config().Get(AppName).(*Framework).Container
}

// ModuleWebservice 每个模块会创建一个新的 webservice，以模块名称作为路由前缀
func ModuleWebservice(obj ioc.ObjectInterface) *restful.WebService {
	modulePath := path.Join("/", obj.Name())

	ws := new(restful.WebService).
		Path(modulePath).
		Consumes("*/*").
		Produces("*/*")

	RootContainer().Add(ws)
	return ws
}
