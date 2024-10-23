package ip2region

import (
	"context"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/qiaogy91/ioc"
	"github.com/qiaogy91/ioc/config/log"
	"log/slog"
)

var (
	ins = &Ip2Region{
		ObjectImpl: ioc.ObjectImpl{},
		FilePath:   "etc/ip2region.xdb",
	}
)

type Ip2Region struct {
	ioc.ObjectImpl
	log      *slog.Logger
	search   *xdb.Searcher
	FilePath string `json:"filePath" yaml:"filePath"`
}

func (i *Ip2Region) Name() string  { return AppName }
func (i *Ip2Region) Priority() int { return 111 }
func (i *Ip2Region) Init() {
	i.log = log.Sub(AppName)

	bs, err := xdb.LoadContentFromFile(i.FilePath)
	if err != nil {
		panic(err)
	}

	i.search, err = xdb.NewWithBuffer(bs)
	if err != nil {
		panic(err)
	}
}
func (i *Ip2Region) Close(ctx context.Context) error {
	i.search.Close()
	return nil
}

func (i *Ip2Region) SearchByStr(ip string) (string, error) {
	return i.search.SearchByStr(ip)
}

func init() {
	ioc.Config().Registry(ins)
}
