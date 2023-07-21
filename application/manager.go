package application

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/save95/go-pkg/framework/logger"
	"github.com/save95/xlog"
)

type manager struct {
	apps []IApplication
	log  xlog.XLog
}

// NewManager 创建 APP 管理器
// opts 支持设置日志、注册应用
// 日志必须实现了 xlog.XLog 接口
// 注册应用必须实现了 IApplication 接口
func NewManager(opts ...interface{}) IManager {
	m := &manager{}

	for _, opt := range opts {
		// 设置 log
		if v, ok := opt.(xlog.XLog); ok {
			m.log = v
			continue
		}

		// 注册应用
		if v, ok := opt.(IApplication); ok {
			m.Register(v)
			continue
		}
	}

	if m.log == nil {
		m.log = logger.NewDefaultLogger()
	}

	return m
}

// Register 注册应用
func (m *manager) Register(app IApplication) {
	m.apps = append(m.apps, app)
}

// Run 启动应用
func (m *manager) Run() {
	// 启动 app
	for i := range m.apps {
		if err := m.apps[i].Start(); err != nil {
			log.Fatalf("app start failed: %+v\n", err)
		}
	}

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	m.log.Info("Shutting down server...")

	// 关闭
	for i := range m.apps {
		if err := m.apps[i].Shutdown(); err != nil {
			m.log.Error("Server forced to shutdown:", err)
			os.Exit(1)
		}
	}

	m.log.Info("Server exiting")
}
