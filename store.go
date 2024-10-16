package ioc

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sort"
	"strings"
)

// Container 容器
type Container []*Namespace

func (c *Container) Len() int {
	return len(*c)
}

func (c *Container) Less(i, j int) bool { return (*c)[i].Priority < (*c)[j].Priority }

func (c *Container) Swap(i, j int) {
	(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
}

func (c *Container) Namespace(ns string) *Namespace {
	for _, item := range *c {
		if item.NsName == ns {
			return item
		}
	}

	// 获取一个名称空间的store 要永远返回正确且有值
	newStore := &Namespace{
		NsName:   ns,
		Priority: 0,
		Items:    []ObjectInterface{},
	}
	*c = append(*c, newStore)
	return newStore
}

func (c *Container) LoadConfig(filePath string) {
	// （1）将文件内容读取出来
	bs, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	// （2）结构化到一个map 中
	cfg := map[string]map[string]any{}
	if err = yaml.Unmarshal(bs, &cfg); err != nil {
		panic(err)
	}

	// （3）为Ioc 中每一个对象应用配置，map -> bs -> struct
	for _, ns := range *c {
		for _, obj := range ns.Items {
			bs, err := json.Marshal(cfg[obj.Name()]) // map  -> bs
			if err != nil {
				panic(err)
			}

			if err := json.Unmarshal(bs, obj); err != nil {
				panic(err)

			} // bs -> struct
		}
	}
}

func (c *Container) Init() {
	sort.Sort(c)
	for _, item := range *c {
		item.Init()
	}
}
func (c *Container) Close(ctx context.Context) error {
	var errs []string
	sort.Sort(sort.Reverse(c))
	for _, item := range *c {
		if err := item.Close(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("close object errs: %s", strings.Join(errs, "|"))
	}
	return nil
}

// Namespace 容器接口实现
type Namespace struct {
	NsName   string
	Priority int
	Items    []ObjectInterface
}

func (s *Namespace) Name() string { return s.NsName }
func (s *Namespace) List() []string {
	var arr []string
	for _, item := range s.Items {
		arr = append(arr, item.Name())
	}
	return arr
}

func (s *Namespace) Get(name string) ObjectInterface {
	for _, item := range s.Items {
		if item.Name() == name {
			return item
		}
	}
	return nil
}

func (s *Namespace) Registry(o ObjectInterface) {
	if item := s.Get(o.Name()); item != nil {
		return
	}
	s.Items = append(s.Items, o)
}

func (s *Namespace) Len() int {
	return len(s.Items)
}

func (s *Namespace) Less(i, j int) bool { return s.Items[i].Priority() < s.Items[j].Priority() }

func (s *Namespace) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

func (s *Namespace) Init() {
	sort.Sort(s)
	for _, item := range s.Items {
		item.Init()
	}
}

func (s *Namespace) Close(ctx context.Context) error {
	var errs []string
	sort.Sort(sort.Reverse(s))
	for _, item := range s.Items {
		if err := item.Close(ctx); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("close object errs: %s", strings.Join(errs, "|"))
	}
	return nil
}
