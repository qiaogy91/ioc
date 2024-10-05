package main

import (
	_ "github.com/qiaogy91/ioc/apps/health/restful"
	_ "github.com/qiaogy91/ioc/apps/metrics/restful"
	_ "github.com/qiaogy91/ioc/example/apps"
	//_ "github.com/qiaogy91/ioc/example/docs"         // gin doc
	//_ "github.com/qiaogy91/ioc/apps/health/gin"  // gin health
	//_ "github.com/qiaogy91/ioc/apps/metrics/gin" // gin metric
	//_ "github.com/qiaogy91/ioc/apps/swagger/gin" // gin swagger
	//_ "github.com/qiaogy91/ioc/config/cors/gin"  // gin cors
	_ "github.com/qiaogy91/ioc/apps/swagger/restful"
	_ "github.com/qiaogy91/ioc/config/cors/restful"
	"github.com/qiaogy91/ioc/server"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io:8080
// @BasePath /v2
func main() {
	server.Start()
}
