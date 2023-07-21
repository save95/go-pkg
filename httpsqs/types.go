package httpsqs

import (
	"context"
	"fmt"
	"time"
)

// IClient HTTPSQS 队列客户端
// @link http://blog.zyan.cc/httpsqs/
type IClient interface {
	// Put 入队列（将文本消息放入队列）
	// 返回：当前队列的读取位置点 pos，及可能存在的错误
	Put(ctx context.Context, name, data string) (pos int64, err error)

	// Get 出队列（从队列中取出文本消息）
	// 返回 文本消息 data，当前队列的读取位置点 pos，及可能存在的错误
	Get(ctx context.Context, name string) (data string, pos int64, err error)

	// Status 查看队列状态
	Status(ctx context.Context, name string) (*Status, error)

	// View 查看指定队列位置点的内容
	// 跟一般的队列系统不同的是，HTTPSQS 可以查看指定队列ID（队列点）的内容，包括未出、已出的队列内容。
	// 可以方便地观测进入队列的内容是否正确。另外，假设有一个发送手机短信的队列，由客户端守护进程从队列
	// 中取出信息，并调用“短信网关接口”发送短信。但是，如果某段时间“短信网关接口”有故障，而这段时间队列
	// 位置点300~900的信息已经出队列，但是发送短信失败，我们还可以在位置点300~900被覆盖前，查看到这些
	// 位置点的内容，作相应的处理。
	View(ctx context.Context, name string, pos int64) (string, error)

	// Reset 重置指定队列
	Reset(ctx context.Context, name string) error

	// SetMaxQueue 更改指定队列的最大队列数量。默认的最大队列长度（100万条）
	SetMaxQueue(ctx context.Context, name string, max int) error

	// SetSyncTime 不停止服务的情况下，修改定时刷新内存缓冲区内容到磁盘的间隔时间
	// 从HTTPSQS 1.3版本开始支持此功能。
	// 默认间隔时间：5秒 或 httpsqs -s <second> 参数设置的值。
	SetSyncTime(ctx context.Context, name string, duration time.Duration) error
}

// IHandler HTTPSQS 队列消费者处理接口
type IHandler interface {
	// QueueName 需要处理的队列名
	QueueName() string

	// GetClient 获取 HTTPSQS 客户端
	GetClient() (IClient, error)

	// OnBefore 前置操作
	OnBefore(ctx context.Context) error

	// Handle 消费队列数据
	Handle(ctx context.Context, data string, pos int64) error

	// OnFailed 失败回调
	OnFailed(ctx context.Context, data string, err error)
}

// Config HTTPSQS 队列客户端配置
type Config struct {
	Addr     string
	Password string
	Timeout  time.Duration
}

// Status 队列状态
type Status struct {
	Name     string `json:"name"`     // 队列名
	MaxQueue int64  `json:"maxqueue"` // 最大队列数量
	PutPos   int64  `json:"putpos"`   // 当前队列的写入位置点
	PutLap   int64  `json:"putlap"`   // 当前队列的写入圈数
	GetPos   int64  `json:"getpos"`   // 当前队列的读取位置点
	GetLap   int64  `json:"getlap"`   // 当前队列的读取圈数
	Unread   int64  `json:"unread"`   // 未读消息数
}

func (s Status) String() string {
	return fmt.Sprintf(`
HTTP Simple Queue Service v1.7
------------------------------
Queue Name: %s
Maximum number of queues: %d
Put position of queue (%dst lap): %d
Get position of queue (%dst lap): %d
Number of unread queue: %d
`, s.Name, s.MaxQueue, s.PutLap, s.PutPos, s.GetLap, s.GetPos, s.Unread)
}
