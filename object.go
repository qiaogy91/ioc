package ioc

import "context"

var _ ObjectInterface = &ObjectImpl{}

type ObjectImpl struct {
}

func (o *ObjectImpl) Init() {}

func (o *ObjectImpl) Name() string {
	return ""
}

func (o *ObjectImpl) Priority() int {
	return 0
}

func (o *ObjectImpl) Close(ctx context.Context) error {
	return nil
}

func (o *ObjectImpl) Meta() ObjectMeta {
	return ObjectMeta{
		PathPrefix: "",
		Extra:      map[string]string{},
	}
}
