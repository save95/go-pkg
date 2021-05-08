package locker

import "time"

// ILocker 锁约定
type ILocker interface {
	// Lock 获得锁
	Lock(key string) error
	// UnLock 释放锁
	UnLock(key string) error
}

// ILockerTimeout 锁超时约定
type ILockerTimeout interface {
	// SetTimeout 设置超时时间
	SetTimeout(expire time.Duration) error
}
