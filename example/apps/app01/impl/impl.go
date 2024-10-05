package impl

import (
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/datasource"
	"github.com/qiaogy91/ioc/config/grpc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/example/apps/app01"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var (
	_ ioc.ObjectInterface = &Impl{} // 符合Ioc 对象约束
	_ app01.Service       = &Impl{} // 符合Grpc 对象约束
)

type Impl struct {
	ioc.ObjectImpl                   // 提供注册到 ioc 的能力
	app01.UnimplementedServiceServer // 提供实现了gRpc 的能力

	log       *zerolog.Logger
	db        *gorm.DB
	KafkaName string `yaml:"kafkaName"`
}

func (i *Impl) Name() string  { return app01.AppName }
func (i *Impl) Priority() int { return 301 }

func (i *Impl) Init() {
	i.log = log.Sub(app01.AppName)
	i.db = datasource.DB()

	// 注册grpc server
	app01.RegisterServiceServer(grpc.Get().Server(), i)
}

func init() {
	ioc.Controller().Registry(&Impl{})
}
