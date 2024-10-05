package impl

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/qiaogy91/ioc/example/apps/app01"
)

func (i *Impl) CreatTable(ctx context.Context) error {
	return i.db.AutoMigrate(&app01.User{})
}

func (i *Impl) Create(ctx context.Context, req *app01.CreateUserReq) (*app01.User, error) {
	if err := validator.New().Struct(req); err != nil {
		return nil, err
	}

	ins := &app01.User{Spec: req}
	if err := i.db.WithContext(ctx).Model(&app01.User{}).Create(ins).Error; err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *Impl) List(ctx context.Context, req *app01.ListUserReq) (*app01.UserSet, error) {
	ins := &app01.UserSet{}
	if err := i.db.WithContext(ctx).Find(&ins.Items).Count(&ins.Total).Error; err != nil {
		return nil, err
	}
	return ins, nil
}
