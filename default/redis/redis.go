package redis

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/redis/go-redis/v9"
	"log/slog"
)

var (
	ins = &Redis{
		Address:  []string{"127.0.0.1:6379"},
		Username: "default",
		Password: "redhat",
		Database: 1,
	}
)

type Redis struct {
	ioc.ObjectImpl
	log      *slog.Logger
	Address  []string `json:"address" yaml:"address"`
	Username string   `json:"username" yaml:"username"`
	Password string   `json:"password" yaml:"password"`
	Database int      `json:"database" yaml:"database"`
	client   redis.UniversalClient
	lock     *redislock.Client
}

func (rds *Redis) Name() string  { return AppName }
func (rds *Redis) Priority() int { return 303 }
func (rds *Redis) Init() {
	rds.log = log.Sub(AppName)
	rds.client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    rds.Address,
		Username: rds.Username,
		Password: rds.Password,
		DB:       rds.Database,
	})

	if err := rds.client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	rds.lock = redislock.New(rds.client)
}

func (rds *Redis) Close(ctx context.Context) error {
	if rds.client == nil {
		return nil
	}
	rds.log.Info("closed completed", slog.String("namespace", ioc.DefaultNamespace))
	return rds.client.Close()
}
func init() {
	ioc.Default().Registry(ins)
}
