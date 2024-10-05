package ioc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
	"os"
	"sort"
	"strings"
)

// Container 容器
type Container struct {
	store []*NamespaceStore
}

func (c *Container) Len() int {
	return len(c.store)
}

func (c *Container) Less(i, j int) bool { return c.store[i].Priority < c.store[j].Priority }

func (c *Container) Swap(i, j int) {
	c.store[i], c.store[j] = c.store[j], c.store[i]
}

func (c *Container) Namespace(ns string) *NamespaceStore {
	for _, item := range c.store {
		if item.Namespace == ns {
			return item
		}
	}

	// 获取一个名称空间的store 要永远返回正确且有值
	newStore := &NamespaceStore{
		Namespace: ns,
		Priority:  0,
		Items:     []ObjectInterface{},
	}
	c.store = append(c.store, newStore)
	return newStore
}

func (c *Container) LoadConfig(req *LoadConfigReq) {
	// 优先加载环境变量
	if req.ConfigEnv.Enabled {
		for i := range c.store {
			item := c.store[i]
			item.LoadFromEnv(req.ConfigEnv.Prefix)
		}
	}

	// 再加载配置文件
	if req.ConfigFile.Enabled {
		for _, item := range c.store {
			item.LoadFromFile(req.ConfigFile.Path)
		}
	}
}

func (c *Container) Init() {
	sort.Sort(c)
	for _, item := range c.store {
		item.Init()
	}
}

// NamespaceStore 容器接口实现
type NamespaceStore struct {
	Namespace string
	Priority  int
	Items     []ObjectInterface
}

func (s *NamespaceStore) List() []string {
	var arr []string
	for _, item := range s.Items {
		arr = append(arr, item.Name())
	}
	return arr
}

func (s *NamespaceStore) Get(name string) ObjectInterface {
	for _, item := range s.Items {
		if item.Name() == name {
			return item
		}
	}
	return nil
}

func (s *NamespaceStore) Registry(o ObjectInterface) {
	if item := s.Get(o.Name()); item != nil {
		return
	}
	s.Items = append(s.Items, o)
}

func (s *NamespaceStore) Len() int {
	return len(s.Items)
}

func (s *NamespaceStore) Less(i, j int) bool { return s.Items[i].Priority() < s.Items[j].Priority() }

func (s *NamespaceStore) Swap(i, j int) {
	s.Items[i], s.Items[j] = s.Items[j], s.Items[i]
}

func (s *NamespaceStore) Init() {
	sort.Sort(s)
	for _, item := range s.Items {
		item.Init()
	}
}

func (s *NamespaceStore) Close(ctx context.Context) error {
	var errs []string
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

// LoadFromFile 从配置文件中加载对象配置
func (s *NamespaceStore) LoadFromFile(filename string) {
	// 临时map
	cfg := map[string]map[string]any{}
	for _, item := range s.Items {
		cfg[item.Name()] = nil
	}

	// 读取配置文件、填充在临时 map
	if bs, err := os.ReadFile(filename); err != nil {
		panic(err)
	} else {
		if err = yaml.Unmarshal(bs, &cfg); err != nil {
			panic(err)
		}
	}

	// 配置每一个对象
	for _, item := range s.Items {
		// 从cfg 中获取map 字串的bs
		bs, err := json.Marshal(cfg[item.Name()])
		if err != nil {
			panic(err)
		}

		// 将这些bs 反序列化给结构体对象
		if err := json.Unmarshal(bs, item); err != nil {
			panic(err)

		}
	}
}

// LoadFromEnv 从环境比那里中加载对象配置，如果参数 prefix 为空，则以对象的大写名称作为前缀 CMDB_ ；否则以 PREFIX_CMDB 作为前缀进行解析，结果赋值给当前对象
func (s *NamespaceStore) LoadFromEnv(prefix string) {
	for _, item := range s.Items {
		prefixList := strings.ToUpper(item.Name()) + "_"
		if prefix != "" {
			prefixList = fmt.Sprintf("%s_%s", strings.ToUpper(prefix), prefixList)
		}
		err := env.Parse(item, env.Options{
			Prefix: prefixList,
		})
		if err != nil {
			panic(err)
		}
	}
}
