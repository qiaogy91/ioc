package apps

import (
	//_ "github.com/qiaogy91/ioc/example/apps/app01/api/restful" // api restful 实现
	_ "github.com/qiaogy91/ioc/example/apps/app01/api/gin" // api gin 实现
	_ "github.com/qiaogy91/ioc/example/apps/app01/impl"
	_ "github.com/qiaogy91/ioc/example/apps/mcpdemo/api/gin" // mcpdemo MCP 工具注册 (gin)
	_ "github.com/qiaogy91/ioc/example/apps/mcpdemo/impl"    // mcpdemo 业务实现
)
