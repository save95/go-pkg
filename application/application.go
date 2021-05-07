package application

// IApplication 应用约定
type IApplication interface {
	Start() error    // 启动
	Shutdown() error // 关闭
}
