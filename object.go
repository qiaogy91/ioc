package ioc

import "context"

var _ ObjectInterface = &ObjectImpl{}

type ObjectImpl struct{}

func (o *ObjectImpl) Name() string  { panic("must define object name") }
func (o *ObjectImpl) Priority() int { panic("must define object priority") }
func (o *ObjectImpl) Init()         {}

func (o *ObjectImpl) Close(ctx context.Context) error {
	return nil
}
func (o *ObjectImpl) Meta() ObjectMeta {
	return ObjectMeta{
		PathPrefix: "",
		Extra:      map[string]string{},
	}
}
