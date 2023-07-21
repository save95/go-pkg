package application

// IApplication 应用约定
type IApplication interface {
	Start() error    // 启动
	Shutdown() error // 关闭
}

type IManager interface {
	// Register 注册应用
	Register(app IApplication)
	// Run 启动应用
	Run()
}
