package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/qiaogy91/ioc"
	"path"
)

const (
	AppName = "gin-framework"
)

func RootRouter() *gin.Engine {
	return ioc.Config().Get(AppName).(*Framework).Engine
}

// ModuleRouter 每个模块会创建一个新的 GroupRoute，以模块名称作为路由前缀
func ModuleRouter(obj ioc.ObjectInterface) gin.IRouter {
	modulePath := path.Join("/", obj.Name())
	return RootRouter().Group(modulePath)
}
