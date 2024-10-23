package ip2region

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
)

const (
	AppName = config.Ip2RegionName
)

func Get() *Ip2Region {
	return ioc.Config().Get(AppName).(*Ip2Region)
}
