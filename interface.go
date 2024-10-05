package ioc

import "context"

// NamespaceInterface 容器接口约束
type NamespaceInterface interface {
	Registry(obj ObjectInterface)    // 注册对象
	Get(name string) ObjectInterface // 获取对象
	List() []string                  // 打印对象名称列表
	LoadFromEnv(prefix string)       // 从环境变量中加载配置
	LoadFromFile(filename string)    // 从配置文件中加载
}

// ObjectInterface 对象接口约束
type ObjectInterface interface {
	Name() string                    // 对象名称
	Priority() int                   // 对象优先级
	Init()                           // 初始化对象
	Close(ctx context.Context) error // 关闭对象
	Meta() ObjectMeta                // 对象元数据描述信息，扩展使用
}

type ObjectMeta struct {
	PathPrefix string
	Extra      map[string]string
}
