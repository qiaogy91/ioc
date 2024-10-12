package gin

import (
	"github.com/gin-contrib/cors"
	"github.com/qiaogy91/ioc"
	iocgin "github.com/qiaogy91/ioc/config/gin"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
	"time"
)

type CORS struct {
	ioc.ObjectImpl
	log            *slog.Logger
	Enabled        bool     `json:"enabled" yaml:"enabled"`
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
	if c.Enabled {
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
		c.log.Info("Gin CORS enabled")
	}
}

func init() {
	ioc.Config().Registry(&CORS{})
}
