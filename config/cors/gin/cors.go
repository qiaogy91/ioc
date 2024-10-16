package gin

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/qiaogy91/ioc"
	iocgin "github.com/qiaogy91/ioc/config/gin"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
	"time"
)

var (
	ins = &CORS{
		AllowedHeaders: []string{"*"},
		AllowOrigins:   []string{"*"},
		AllowedMethods: []string{"*"},
		ExposeHeaders:  []string{"*"},
		AllowCookies:   true,
		MaxAge:         43200,
	}
)

type CORS struct {
	ioc.ObjectImpl
	log            *slog.Logger
	AllowedHeaders []string `json:"allowedHeaders" yaml:"allowedHeaders"`
	AllowOrigins   []string `json:"allowOrigins" yaml:"allowOrigins"`
	AllowedMethods []string `json:"allowedMethods" yaml:"allowedMethods"`
	ExposeHeaders  []string `json:"exposeHeaders" yaml:"exposeHeaders"`
	AllowCookies   bool     `json:"allowCookies" yaml:"allowCookies"`
	MaxAge         int      `json:"maxAge" yaml:"maxAge"`
}

func (c *CORS) Name() string { return AppName }
func (c *CORS) Priority() int {
	return 106
}
func (c *CORS) Init() {
	c.log = log.Sub(AppName)

	// 将中间件添加到Router中
	r := iocgin.RootRouter() // 将中间件加载到Root 根路由，而非模块的Group 分组路由上
	r.Use(cors.New(cors.Config{
		AllowOrigins:     c.AllowOrigins,
		AllowMethods:     c.AllowedMethods,
		AllowHeaders:     c.AllowedHeaders,
		ExposeHeaders:    c.ExposeHeaders,
		AllowCredentials: c.AllowCookies,
		MaxAge:           time.Duration(c.MaxAge) * time.Second,
		AllowWildcard:    true,
	}))
	c.log.Debug("Gin CORS enabled")
}
func (c *CORS) Close(ctx context.Context) error {
	c.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func init() {
	ioc.Config().Registry(ins)
}
