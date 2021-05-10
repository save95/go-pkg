package dbutil

import (
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type dbFunc func() (*gorm.DB, error)

// 利用该结构减少并发锁竞争
var dbRelation sync.Map

// 获取db
func connectWithoutCache(option *Option) (*gorm.DB, error) {
	config := option.Config
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

	db.LogMode(config.LogMode)
	// 重新设置日志
	if option.Logger != nil {
		db.SetLogger(convertLogger(option.Logger))
	}
	return db, nil
}

// Connect 获取db
func Connect(option *Option) (*gorm.DB, error) {
	if option.Name == "" {
		return nil, errors.New("the db config name invalid")
	}

	var (
		db  *gorm.DB
		err error

		// 用于只初始化一次
		wg sync.WaitGroup
	)
	wg.Add(1)
	fi, loaded := dbRelation.LoadOrStore(option.Name, dbFunc(func() (*gorm.DB, error) {
		// 阻塞直到初始化完成
		wg.Wait()
		return db, err
	}))

	// 已经存在，则直接调用即可
	if loaded {
		return fi.(dbFunc)()
	}

	// 未找到则需要初始化
	db, err = connectWithoutCache(option)

	// 真实的返回db函数，wg释放后
	f := dbFunc(func() (*gorm.DB, error) {
		return db, err
	})

	wg.Done()
	// 重置函数
	dbRelation.Store(option.Name, f)
	return db, err
}
