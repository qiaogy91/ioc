package redis

import (
	"github.com/bsm/redislock"
	"github.com/qiaogy91/ioc"
	"github.com/redis/go-redis/v9"
)

const AppName = "redis"

func GetClient() redis.UniversalClient {
	return Get().client
}

func GetLock() *redislock.Client {
	return Get().lock
}

func Get() *Redis {
	return ioc.Default().Get(AppName).(*Redis)
}
