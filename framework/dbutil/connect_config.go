package dbutil

// ConnectConfig 数据库链接配置
type ConnectConfig struct {
	Dsn         string // 连接
	Driver      string // 数据库类型
	MaxIdle     int    // 最大空闲连接数
	MaxOpen     int    // 最大连接数
	LogMode     bool   // 是否打印SQL
	MaxLifeTime int    // 连接存活时间
}
