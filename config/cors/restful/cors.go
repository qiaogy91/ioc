package restful

import (
	"context"
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

type CORS struct {
	ioc.ObjectImpl
	log            *slog.Logger
	Enabled        bool     `json:"enabled" yaml:"enabled"`
	AllowedHeaders []string `json:"allowedHeaders" yaml:"allowedHeaders"`
	AllowedDomains []string `json:"allowedDomains" yaml:"allowedDomains"`
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
	c.log = log.Sub("cors")

	// 将中间件添加到Router中
	if c.Enabled {
		container := gorestful.RootContainer()
		cors := restful.CrossOriginResourceSharing{
			// 可以的
			ExposeHeaders:  c.ExposeHeaders,
			AllowedHeaders: c.AllowedHeaders,
			AllowedMethods: c.AllowedMethods,
			AllowedDomains: c.AllowedDomains,
			CookiesAllowed: false,
			Container:      container,
			MaxAge:         c.MaxAge,
		}
		container.Filter(cors.Filter)
		c.log.Debug("Restful CORS enabled")
	}
}

func (c *CORS) Close(ctx context.Context) error {
	c.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return nil
}

func init() {
	ioc.Config().Registry(&CORS{})
}
