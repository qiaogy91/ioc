package ioc

import "sync"

const (
	ConfigNamespace     = "configs"
	ControllerNamespace = "controllers"
	DefaultNamespace    = "default"
	ApiNamespace        = "apis"
)

// 定义各个名称空间容器，并定义获取这些容器的方法
var (
	container *Container
	lock      sync.Mutex
)

// GetContainer 获取整个Ioc 容器
func GetContainer() *Container {
	return container
}

func Controller() *Namespace { return container.Namespace(ControllerNamespace) } // 获取 Controller 名称空间
func Config() *Namespace     { return container.Namespace(ConfigNamespace) }     // 获取 Config 名称空间
func Api() *Namespace        { return container.Namespace(ApiNamespace) }        // 获取 Api 名称空间
func Default() *Namespace    { return container.Namespace(DefaultNamespace) }    // 获取 Default 名称空间

// ConfigIocObject 为容器中的对象进行二次配置
func ConfigIocObject(confPath string) error {
	lock.Lock()
	defer lock.Unlock()

	container.LoadConfig(confPath) // 加载配置文件到内存
	container.Init()               // 执行每个对象的 Init() 方法进行初始化
	return nil
}

func init() {
	container = &Container{
		{NsName: ConfigNamespace, Priority: 1, Items: make([]ObjectInterface, 0)},
		{NsName: DefaultNamespace, Priority: 2, Items: make([]ObjectInterface, 0)},
		{NsName: ControllerNamespace, Priority: 3, Items: make([]ObjectInterface, 0)},
		{NsName: ApiNamespace, Priority: 4, Items: make([]ObjectInterface, 0)},
	}
}
