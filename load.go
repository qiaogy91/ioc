package ioc

const (
	ConfigNamespace     = "configs"
	ControllerNamespace = "controllers"
	DefaultNamespace    = "default"
	ApiNamespace        = "apis"
)

// 定义各个名称空间容器，并定义获取这些容器的方法
var (
	isLoaded     bool
	ConfFilePath string
	store        = &Container{
		store: []*NamespaceStore{
			{Namespace: ConfigNamespace, Priority: 1, Items: make([]ObjectInterface, 0)},
			{Namespace: DefaultNamespace, Priority: 2, Items: make([]ObjectInterface, 0)},
			{Namespace: ControllerNamespace, Priority: 3, Items: make([]ObjectInterface, 0)},
			{Namespace: ApiNamespace, Priority: 4, Items: make([]ObjectInterface, 0)},
		},
	}
)

func Controller() *NamespaceStore { return store.Namespace(ControllerNamespace) }
func Config() *NamespaceStore     { return store.Namespace(ConfigNamespace) }
func Api() *NamespaceStore        { return store.Namespace(ApiNamespace) }
func Default() *NamespaceStore    { return store.Namespace(DefaultNamespace) }

// ConfigIocObject 为容器中的对象进行二次配置
func ConfigIocObject(confPath string) error {
	if !isLoaded {
		// 加载配置
		store.LoadConfig(confPath)
		// 初始化对象
		store.Init()
	}

	isLoaded = true
	return nil
}
