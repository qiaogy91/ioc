package restful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/gorestful"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/rs/zerolog"
)

type CORS struct {
	ioc.ObjectImpl
	log            *zerolog.Logger
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
		c.log.Info().Msg("restful cors enabled")
	}
}

func init() {
	ioc.Config().Registry(&CORS{})
}
