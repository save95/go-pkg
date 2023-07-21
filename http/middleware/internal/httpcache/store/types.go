package store

import (
	"errors"
	"net/http"
	"time"
)

var ErrorCacheMiss = errors.New("cache miss error")

type ICacheStore interface {
	// Get 获取缓存，如果未获取到，返回 ErrorCacheMiss 错误
	Get(key string, value *CachedResponse) error

	// Set 设置缓存，如果存在，则覆盖
	Set(key string, value *CachedResponse, expire time.Duration) error

	// Delete 删除缓存，如果不存在，则不处理
	Delete(key string) error
}

// CachedResponse 缓存的响应
type CachedResponse struct {
	Status int
	Header http.Header

	Data []byte
}
