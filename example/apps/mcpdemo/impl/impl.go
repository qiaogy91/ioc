package impl

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"github.com/qiaogy91/ioc/example/apps/mcpdemo"
)

var (
	_ ioc.ObjectInterface = &Impl{}
	_ mcpdemo.Controller  = &Impl{}
)

type Impl struct {
	ioc.ObjectImpl
	log *slog.Logger
}

func (i *Impl) Name() string  { return mcpdemo.ControllerName }
func (i *Impl) Priority() int { return 301 }

func (i *Impl) Init() {
	i.log = log.Sub(mcpdemo.ControllerName)
	i.log.Debug("mcpdemo service initialized")
}

// GetUser 根据 ID 获取用户信息
func (i *Impl) GetUser(ctx context.Context, id int) (*mcpdemo.User, error) {
	i.log.Info("GetUser called", slog.Int("id", id))
	// 模拟数据库查询
	return &mcpdemo.User{
		ID:   id,
		Name: fmt.Sprintf("User_%d", id),
		Age:  20 + id%30,
	}, nil
}

func init() {
	ioc.Controller().Registry(&Impl{})
}
