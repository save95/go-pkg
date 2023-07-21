package httpcache

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/http/middleware/internal/httpcache/store"
	"golang.org/x/sync/singleflight"
)

var sf singleflight.Group

// strategy 缓存处理策略
type strategy struct {
	NeedCached    bool          // 是否需要缓存
	CacheKey      string        // 缓存 Key
	CacheDuration time.Duration // 缓存时效
}

type ruleItem struct {
	withToken  bool
	fields     map[string]struct{} // 用于计算缓存的字段 key。会覆盖 globalSkipFields 规则
	skipFields map[string]struct{} // 不用于计算缓存的 key
	duration   time.Duration       // 缓存时长。会覆盖 globalDuration 规则
}

func (m *ruleItem) String() string {
	bs, _ := json.Marshal(map[string]interface{}{
		"withToken":  m.withToken,
		"fields":     m.fields,
		"skipFields": m.skipFields,
		"duration":   m.duration,
	})
	return string(bs)
}

func newCachedResponse(writer *responseWriter) *store.CachedResponse {
	return &store.CachedResponse{
		Status: writer.Status(),
		Header: writer.Header().Clone(),
		Data:   writer.body.Bytes(),
	}
}

type responseWriter struct {
	gin.ResponseWriter

	body *bytes.Buffer
}

func newResponseWriter(w gin.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (r *responseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r *responseWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

//func (r *responseWriter) Body() *bytes.Buffer {
//	var bs bytes.Buffer
//
//	_, _ = bs.ReadFrom(r.body)
//
//	return &bs
//}
