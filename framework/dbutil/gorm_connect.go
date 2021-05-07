package dbutil

import (
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type (
	dbLog struct{}

	// 重置db日志打印
	dbFunc func() (*gorm.DB, error)
)

func (l *dbLog) Print(v ...interface{}) {
	log.Println(gorm.LogFormatter(v...)...)
}

// 利用该结构减少并发锁竞争
var dbRelation sync.Map

// 获取db
func connectWithoutCache(config *ConnectConfig) (*gorm.DB, error) {
	db, err := gorm.Open(config.Driver, config.Dsn)
	if err != nil {
		return nil, errors.Wrap(err, "db连接异常")
	}

	if config.MaxIdle > 0 {
		db.DB().SetMaxIdleConns(config.MaxIdle)
	}
	if config.MaxOpen > 0 {
		db.DB().SetMaxOpenConns(config.MaxOpen)
	}
	if config.MaxLifeTime > 0 {
		db.DB().SetConnMaxLifetime(time.Duration(config.MaxLifeTime) * time.Second)
	}

	// 重新设置日志
	db.SetLogger(&dbLog{})
	db.LogMode(config.LogMode)
	return db, nil
}

// Connect 获取db
func Connect(name string, config *ConnectConfig) (*gorm.DB, error) {
	if name == "" {
		return nil, errors.New("the db config name invalid")
	}

	var (
		db  *gorm.DB
		err error

		// 用于只初始化一次
		wg sync.WaitGroup
	)
	wg.Add(1)
	fi, loaded := dbRelation.LoadOrStore(name, dbFunc(func() (*gorm.DB, error) {
		// 阻塞直到初始化完成
		wg.Wait()
		return db, err
	}))

	// 已经存在，则直接调用即可
	if loaded {
		return fi.(dbFunc)()
	}

	// 未找到则需要初始化
	db, err = connectWithoutCache(config)

	// 真实的返回db函数，wg释放后
	f := dbFunc(func() (*gorm.DB, error) {
		return db, err
	})

	wg.Done()
	// 重置函数
	dbRelation.Store(name, f)
	return db, err
}
