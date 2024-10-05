package datasource

import (
	"github.com/qiaogy91/ioc"
	"gorm.io/gorm"
)

const (
	AppName = "datasource"
)

func Get() *DataSource {
	return ioc.Config().Get(AppName).(*DataSource)
}

func DB() *gorm.DB {
	return Get().db
}
