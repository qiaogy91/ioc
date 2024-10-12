package datasource

import (
	"context"
	"fmt"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/config/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"log/slog"
)

var _ ioc.ObjectInterface = &DataSource{}

type DataSource struct {
	ioc.ObjectImpl
	Trace    bool   `json:"trace" yaml:"trace"`
	Host     string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	Port     int    `json:"port" yaml:"port" toml:"port" env:"PORT"`
	DB       string `json:"database" yaml:"database" toml:"database" env:"DB"`
	Username string `json:"username" yaml:"username" toml:"username" env:"USERNAME"`
	Password string `json:"password" yaml:"password" toml:"password" env:"PASSWORD"`
	Debug    bool   `json:"debug" yaml:"debug" toml:"debug" env:"DEBUG"`

	db  *gorm.DB
	log *slog.Logger
}

func (ds *DataSource) Name() string {
	return AppName
}

func (ds *DataSource) Priority() int {
	return 105
}

func (ds *DataSource) Init() {
	ds.log = log.Sub(AppName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ds.Username,
		ds.Password,
		ds.Host,
		ds.Port,
		ds.DB,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true, // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
		SkipDefaultTransaction: true, // 对于写操作，默认Gorm 为了数据的完整性将其封装在事务中运行。如果没有这方面要求可关闭，性能会提升30%
	})
	if err != nil {
		panic(err)
	}

	if ds.Debug {
		db = db.Debug()
	}

	ds.db = db

	// 开启Trace
	if trace.Get().Enable && ds.Trace {
		if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
			panic(err)
		}
		ds.log.Info("mysql trace enabled")
	}
}

func (ds *DataSource) Close(ctx context.Context) error {
	if ds.db == nil {
		return nil
	}

	d, err := ds.db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}

func init() {
	ioc.Config().Registry(&DataSource{})
}
