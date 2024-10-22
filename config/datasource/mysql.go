package datasource

import (
	"context"
	"fmt"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"log/slog"
)

var (
	_   ioc.ObjectInterface = &DataSource{}
	ins                     = &DataSource{
		Otlp:     false,
		Host:     "127.0.0.1",
		Port:     3306,
		DB:       "must_set",
		Username: "root",
		Password: "redhat",
		Debug:    true,
	}
)

type DataSource struct {
	ioc.ObjectImpl
	Otlp     bool   `json:"otlp" yaml:"otlp"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	DB       string `json:"database" yaml:"database"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
	Debug    bool   `json:"debug" yaml:"debug"`

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

	// 开启 PrepareStmt 则会在执行 SQL 时缓存预编译语句。
	// GORM 在启用该功能时，会在 ConnPool 中使用一个自定义类型来支持预编译语句的缓存机制，即导致 db.ConnPool 不再是 *sql.DB 类型
	// go-gorm/opentelemetry 会通过判断 db.ConnPool 是否是 *sql.DB 类型，来决定是否注册DB 相关的Metrics
	// 因此开启该参数后，会导致无法收集 DB 相关的Metrics，因此当遥测开关打开时，需要关闭这个Gorm 参数
	conf := &gorm.Config{SkipDefaultTransaction: true}
	if !ds.Otlp {
		conf.PrepareStmt = true // 执行任何 SQL 时都创建并缓存预编译语句，可以提高后续的调用速度
	}

	db, err := gorm.Open(mysql.Open(dsn), conf)
	if err != nil {
		panic(err)
	}

	if ds.Debug {
		db = db.Debug()
	}
	ds.db = db

	// 开启Trace
	if ds.Otlp {
		ioc.OtlpMustEnabled()
		if err := db.Use(tracing.NewPlugin()); err != nil {
			panic(err)
		}
		ds.log.Debug("mysql Otlp enabled")
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
	ds.log.Debug("closed completed", slog.String("namespace", ioc.ConfigNamespace))
	return d.Close()
}

func init() {
	ioc.Config().Registry(ins)
}
