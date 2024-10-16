package datasource

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config"
	"gorm.io/gorm"
)

const (
	AppName = config.DatasourceName
)

func Get() *DataSource {
	return ioc.Config().Get(AppName).(*DataSource)
}

func DB() *gorm.DB {
	return Get().db
}
