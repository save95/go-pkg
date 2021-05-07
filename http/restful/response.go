package restful

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/save95/go-pkg/constant"
	"github.com/save95/go-pkg/http/types"
	"github.com/save95/go-pkg/utils/strutil"
	"github.com/save95/xerror"
	"github.com/save95/xlog"
)

type response struct {
	ctx    *gin.Context
	logger xlog.XLog
}

// NewResponse 创建 Restful 标准响应生成器
func NewResponse(ctx *gin.Context, logger xlog.XLog) *response {
	return &response{ctx: ctx, logger: logger}
}

// SetHeader 设置请求头
func (r *response) SetHeader(key, value string) *response {
	// 必须使用自定义头 X- 开始才设置，否则跳过
	if !strings.HasPrefix(key, "X-") && !strings.HasPrefix(key, "x-") {
		return r
	}

	r.ctx.Header(key, value)
	return r
}

// Retrieve 查询单个资源的响应
func (r *response) Retrieve(entity interface{}) {
	//r.ctx.Header("Content-MD5", fmt.Sprintf("%x", md5.Sum([]byte())))
	if entity == nil {
		r.ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	r.ctx.AbortWithStatusJSON(http.StatusOK, entity)
}

// ListWithPagination 分页列表的响应
func (r *response) ListWithPagination(totalRow uint, entities interface{}) {
	tk := reflect.TypeOf(entities).Kind()
	if tk != reflect.Slice && tk != reflect.Array {
		r.WithError(xerror.New("response data type error"))
		return
	}

	// 设置总记录数
	r.ctx.Header("X-Total-Count", strconv.Itoa(int(totalRow)))

	// 解析URL，Query string
	currentUri := r.ctx.Request.RequestURI
	urls, err := url.Parse(currentUri)
	if nil != err {
		r.WithError(xerror.Wrap(err, "parse uri failed"))
		return
	}

	qs := urls.Query()
	start := strutil.ToInt(qs.Get("start"))
	limit := strutil.ToIntWith(qs.Get("limit"), constant.DefaultRequestLimit)

	// 计算分页信息
	page := uint(math.Ceil(float64(start/limit)) + 1)
	count := uint(math.Max(1, float64(totalRow/uint(limit))))

	// 设置分页信息
	r.ctx.Header("X-Pagination-Info", fmt.Sprintf(
		`count="%d", rows="%d", current="%d", size="%d"`,
		count,
		totalRow,
		page,
		limit,
	))

	// 计算分页url
	firstUri := r.ComputePaginateUri(urls, 0)

	prevStart := int(math.Max(0, float64(start-limit)))
	prevUri := r.ComputePaginateUri(urls, prevStart)

	nextStart := int(math.Min(float64((count-1)*uint(limit)), float64(page*uint(limit))))
	nextUri := r.ComputePaginateUri(urls, nextStart)

	lastStart := int(math.Max(0, float64((count-1)*uint(limit))))
	lastUri := r.ComputePaginateUri(urls, lastStart)

	links := fmt.Sprintf(
		`<%s>; rel="self", <%s>; rel="previous", <%s>; rel="next", <%s>; rel="first", <%s>; rel="last"`,
		currentUri,
		prevUri,
		nextUri,
		firstUri,
		lastUri,
	)
	r.ctx.Header("Link", links)

	//r.ctx.Header("Content-MD5", "")

	if reflect.ValueOf(entities).IsNil() {
		entities = make([]interface{}, 0)
	}
	r.ctx.AbortWithStatusJSON(http.StatusOK, entities)
}

func (r *response) ComputePaginateUri(urls *url.URL, start int) string {
	qs := urls.Query()
	qs.Set("start", strconv.Itoa(start))
	if start == 0 {
		qs.Del("start")
	}

	if len(qs.Encode()) == 0 {
		return urls.Path
	}

	return fmt.Sprintf("%s?%s", urls.Path, qs.Encode())
}

// ListWithMoreFlag 查询列表的响应
func (r *response) ListWithMoreFlag(hasMore bool, entities interface{}) {
	tk := reflect.TypeOf(entities).Kind()
	if tk != reflect.Slice && tk != reflect.Array {
		r.WithError(xerror.New("response data type error"))
		return
	}

	//if len(entities) == 0 {
	//	hasMore = false
	//}

	r.ctx.Header("X-More-Resource", strconv.FormatBool(hasMore))

	if reflect.ValueOf(entities).IsNil() {
		entities = make([]interface{}, 0)
	}
	r.ctx.AbortWithStatusJSON(http.StatusOK, entities)
}

// Post 新增请求的响应
func (r *response) Post(entity interface{}) {
	if nil == entity {
		r.WithError(xerror.New("post must has response entity"))
		return
	}

	r.ctx.AbortWithStatusJSON(http.StatusCreated, entity)
}

// Put 全量更新资源的响应
func (r *response) Put(entity interface{}) {
	if nil == entity {
		r.WithError(xerror.New("put must has response entity"))
		return
	}

	r.ctx.AbortWithStatusJSON(http.StatusCreated, entity)
}

// Patch 部分更新资源的响应
// 部分 cdn 服务商不支持 http patch 方法，如 阿里云
func (r *response) Patch(entity interface{}) {
	if nil == entity {
		r.ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	r.ctx.AbortWithStatusJSON(http.StatusCreated, entity)
}

// Delete 删除的响应
func (r *response) Delete(err error) {
	if nil != err {
		r.WithError(err)
		return
	}

	r.ctx.AbortWithStatus(http.StatusNoContent)
}

// WithMessage 通过 json 响应文本消息: {"message": "something..."}
func (r *response) WithMessage(msg string) {
	if len(msg) == 0 {
		msg = "success"
	}

	r.ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
		"message": msg,
	})
}

// WithBody 响应文本消息
func (r *response) WithBody(body string) {
	r.ctx.String(http.StatusOK, "%s", body)
}

// WithError 响应错误消息(HttpStatus!=200)
func (r *response) WithError(err error) {
	rq := r.ctx.Request
	if stx, se := types.ParserHttpContext(r.ctx); nil == se {
		bs := stx.Value(constant.HttpCustomRawRequestBodyKey).([]byte)
		rq.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
	}

	r.logger.Error(newRequestError(rq, err))

	if e, ok := err.(xerror.XError); ok {
		r.ctx.AbortWithStatusJSON(e.HttpStatus(), gin.H{
			"message": e.String(),
		})
		return
	}

	r.ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"message": err.Error(),
	})
}
