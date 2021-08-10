package dbutil

import (
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type dbFunc func() (*gorm.DB, error)

// 利用该结构减少并发锁竞争
var dbRelation sync.Map

// 获取db
func connectWithoutCache(dialect gorm.Dialector, option *Option) (*gorm.DB, error) {
	// 连接 db
	db, err := gorm.Open(dialect, &gorm.Config{
		Logger: newLogger(option.Logger),
	})
	if err != nil {
		return nil, errors.Wrap(err, "db连接异常")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "db连接异常")
	}

	config := option.Config
	if config.MaxIdle > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdle)
	}
	if config.MaxOpen > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpen)
	}
	if config.MaxLifeTime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.MaxLifeTime) * time.Second)
	}

	return db, nil
}

// Connect 获取db
func Connect(option *Option) (*gorm.DB, error) {
	if option == nil {
		return nil, errors.New("db connect option empty")
	}
	if option.Name == "" {
		return nil, errors.New("the db config name invalid")
	}
	if option.Config == nil {
		return nil, errors.New("db config empty")
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

	// 配置转成方言
	dialect, err := toDialect(option.Config.Driver, option.Config.Dsn)
	if nil != err {
		return nil, err
	}

	// 未找到则需要初始化
	db, err = connectWithoutCache(dialect, option)

	// 真实的返回db函数，wg释放后
	f := dbFunc(func() (*gorm.DB, error) {
		return db, err
	})

	wg.Done()
	// 重置函数
	dbRelation.Store(option.Name, f)
	return db, err
}

// ConnectWith 通过方言获取db
func ConnectWith(dialect gorm.Dialector, option *Option) (*gorm.DB, error) {
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
	db, err = connectWithoutCache(dialect, option)

	// 真实的返回db函数，wg释放后
	f := dbFunc(func() (*gorm.DB, error) {
		return db, err
	})

	wg.Done()
	// 重置函数
	dbRelation.Store(option.Name, f)
	return db, err
}
