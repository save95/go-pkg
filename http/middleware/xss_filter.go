package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/save95/xerror"
)

type xssFilter struct {
	xssRuleItem

	debug bool

	// 路由特殊规则
	routePolicies map[string]*xssRuleItem
	// 直接跳过的路由
	skipRoutes map[string]struct{}
}

type xssRuleItem struct {
	policy    *bluemonday.Policy
	skipField map[string]struct{}
}

// XSSPolicy XSS 策略
type XSSPolicy uint8

const (
	// XSSPolicyNone 无过滤
	XSSPolicyNone XSSPolicy = iota
	// XSSPolicyStrict 过滤所有HTML元素及其属性
	XSSPolicyStrict
	// XSSPolicyUGC 过滤不安全的HTML元素和属性，如：iframes, object, embed, styles, script
	XSSPolicyUGC
)

// XSSGlobalPolicy 指定全局过滤策略
func XSSGlobalPolicy(p XSSPolicy) func(xf *xssFilter) {
	return func(xf *xssFilter) {
		xf.policy = xf.makePolicy(p)
	}
}

// XSSDebug 设置调试模式
func XSSDebug() func(xf *xssFilter) {
	return func(xf *xssFilter) {
		xf.debug = true
	}
}

// XSSGlobalSkipFields 指定全局忽略字段
func XSSGlobalSkipFields(fields ...string) func(xf *xssFilter) {
	return func(xf *xssFilter) {
		xf.skipField = xf.makeSkipFields(fields)
	}
}

// XSSRoutePolicy 指定路由策略
// routeRule 路由规则，如果路由包含该字符串则匹配成功
func XSSRoutePolicy(routeRule string, policy XSSPolicy, skipFields ...string) func(xf *xssFilter) {
	return func(xf *xssFilter) {
		if policy == XSSPolicyNone {
			xf.skipRoutes[routeRule] = struct{}{}
			return
		}

		xf.routePolicies[routeRule] = &xssRuleItem{
			policy:    xf.makePolicy(policy),
			skipField: xf.makeSkipFields(skipFields),
		}
	}
}

// XSSFilter XSS 过滤
// usage:
// r.Use(middleware.XSSFilter(
//  	// middleware.XSSDebug(),
//  	middleware.XSSGlobalPolicy(middleware.XSSPolicyStrict),
//  	middleware.XSSGlobalSkipFields("password"),
//  	middleware.XSSRoutePolicy("admin", middleware.XSSPolicyUGC),
//  	middleware.XSSRoutePolicy("/callback/", middleware.XSSPolicyNone),
//  	middleware.XSSRoutePolicy("/endpoint", middleware.XSSPolicyNone),
//  	middleware.XSSRoutePolicy("/ping", middleware.XSSPolicyNone),
// ))
func XSSFilter(opts ...func(xf *xssFilter)) gin.HandlerFunc {
	xf := &xssFilter{
		xssRuleItem: xssRuleItem{
			skipField: make(map[string]struct{}, 0),
		},
		routePolicies: make(map[string]*xssRuleItem, 0),
		skipRoutes:    make(map[string]struct{}, 0),
	}
	//xf.policy = xf.makePolicy(XSSPolicyStrict)

	for _, opt := range opts {
		opt(xf)
	}

	return xf.filter()
}

func (xf *xssFilter) makePolicy(p XSSPolicy) *bluemonday.Policy {
	switch p {
	case XSSPolicyNone:
		return nil
	case XSSPolicyStrict:
		return bluemonday.StrictPolicy()
	case XSSPolicyUGC:
		return bluemonday.UGCPolicy()
	default:
		return bluemonday.StrictPolicy()
	}
}

func (xf *xssFilter) makeSkipFields(fields []string) map[string]struct{} {
	vals := make(map[string]struct{}, 0)

	for _, field := range fields {
		vals[field] = struct{}{}
	}

	return vals
}

func (xf *xssFilter) filter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 指定了忽略路由，直接跳过
		for u := range xf.skipRoutes {
			if strings.Contains(ctx.FullPath(), u) {
				xf.debugf("xss filter hit skip route, skip\n")
				ctx.Next()
				return
			}
		}

		// 未指定全局规则，而且未指定路由规则，直接跳过
		if xf.policy == nil {
			skip := true
			for u := range xf.routePolicies {
				if strings.Contains(ctx.FullPath(), u) {
					skip = false
					break
				}
			}
			if skip {
				xf.debugf("xss filter no global policy, not hit route rule, skip\n")
				ctx.Next()
				return
			}
		}

		var err error

		switch ctx.Request.Method {
		case http.MethodGet:
			err = xf.filterQueryString(ctx)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := ctx.Request.Header.Get("Content-Type")
			switch ct {
			case "application/json":
				err = xf.filterJSON(ctx)
			case "application/x-www-form-urlencoded":
				err = xf.filterFormData(ctx)
			default:
				if strings.Contains(ct, "multipart/form-data") {
					err = xf.filterMultiPartFormData(ctx)
				}
			}
		}

		if nil != err {
			if xe, ok := err.(xerror.XError); ok {
				err = xe.Unwrap()
			}
			log.Printf("xss filter err: %+v\n", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.Next()
	}
}

func (xf *xssFilter) filterXSS(fullPath, key, val string) string {
	for route, item := range xf.routePolicies {
		if strings.Contains(fullPath, route) {
			if _, ok := item.skipField[key]; ok {
				xf.debugf("xss filter hit route skip field, return origin value\n")
				return val
			}
			if item.policy == nil {
				xf.debugf("xss filter hit route rule, none policy, return origin value\n")
				return val
			}
			xf.debugf("xss filter hit route rule, return sanitize value\n")
			return item.policy.Sanitize(val)
		}
	}

	if _, ok := xf.skipField[key]; ok {
		xf.debugf("xss filter hit global skip field, return origin value\n")
		return val
	}

	if xf.policy == nil {
		xf.debugf("xss filter hit global policy, return sanitize value\n")
		return val
	}

	xf.debugf("xss filter hit global policy, return sanitize value\n")
	return xf.policy.Sanitize(val)
}

func (xf *xssFilter) debugf(format string, vals ...interface{}) {
	if xf.debug {
		log.Printf(format, vals...)
	}
}

func (xf *xssFilter) filterQueryString(ctx *gin.Context) error {
	params := ctx.Request.URL.Query()
	xf.debugf("xss filter input query string: %s\n", params.Encode())
	for key, items := range params {
		params.Del(key)
		for _, val := range items {
			val = xf.filterXSS(ctx.FullPath(), key, val)
			if params.Has(key) {
				params.Add(key, val)
			} else {
				params.Set(key, val)
			}
		}
	}

	xf.debugf("xss filter output query string: %s\n", params.Encode())
	ctx.Request.URL.RawQuery = params.Encode()
	return nil
}

func (xf *xssFilter) filterJSON(ctx *gin.Context) error {
	body := ctx.Request.Body
	if body == nil || body == http.NoBody {
		return nil
	}

	d := json.NewDecoder(body)
	d.UseNumber()

	var val interface{}
	if err := d.Decode(&val); nil != err {
		return xerror.Wrap(err, "json decode failed")
	}

	xf.debugf("xss filter input json: %s\n", val)

	var data interface{}
	switch val.(type) {
	case map[string]interface{}:
		vals := make(map[string]interface{}, 0)
		for k, v := range val.(map[string]interface{}) {
			vals[k] = xf.filterJsonValue(ctx.FullPath(), k, v)
		}
		data = vals
	case []interface{}:
		vals := make([]interface{}, 0)
		for _, v := range val.([]interface{}) {
			vals = append(vals, xf.filterJsonValue(ctx.FullPath(), "", v))
		}
		data = vals
	default:
		data = val
	}

	var bf bytes.Buffer
	encode := json.NewEncoder(&bf)
	encode.SetEscapeHTML(false)
	if err := encode.Encode(data); nil != err {
		return xerror.Wrap(err, "json encode failed")
	}

	xf.debugf("xss filter output json: %s\n", bf.String())
	ctx.Request.Body = ioutil.NopCloser(&bf)
	return nil
}

func (xf *xssFilter) filterJsonValue(fullPath, key string, val interface{}) interface{} {
	switch val.(type) {
	case map[string]interface{}:
		vals := make(map[string]interface{}, 0)
		for k, v := range val.(map[string]interface{}) {
			vals[k] = xf.filterJsonValue(fullPath, key, v)
		}
		return vals
	case []interface{}:
		vals := make([]interface{}, 0)
		for _, v := range val.([]interface{}) {
			vals = append(vals, xf.filterJsonValue(fullPath, key, v))
		}
		return vals
	case string:
		return xf.filterXSS(fullPath, key, val.(string))
	default:
		return val
	}
}

func (xf *xssFilter) filterFormData(ctx *gin.Context) error {
	body := ctx.Request.Body
	if body == nil || body == http.NoBody {
		return nil
	}

	// https://golang.org/src/net/http/httputil/dump.go
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(body); err != nil {
		return xerror.Wrap(err, "read from failed")
	}

	xf.debugf("xss filter input x-form-data: %s\n", buf.String())

	m, err := url.ParseQuery(buf.String())
	if err != nil {
		return xerror.Wrap(err, "parse query failed")
	}

	var bf bytes.Buffer
	for key, v := range m {
		val := url.QueryEscape(v[0])
		if _, ok := xf.skipField[key]; !ok {
			val = url.QueryEscape(xf.filterXSS(ctx.FullPath(), key, v[0]))
		}

		if bf.Len() > 0 {
			bf.WriteByte('&')
		}
		bf.WriteString(key)
		bf.WriteByte('=')
		bf.WriteString(val)
	}

	xf.debugf("xss filter output x-form-data: %s\n", bf.String())
	ctx.Request.Body = ioutil.NopCloser(&bf)
	return nil
}

func (xf *xssFilter) filterMultiPartFormData(ctx *gin.Context) error {
	body := ctx.Request.Body
	if body == nil || body == http.NoBody {
		return nil
	}

	ct := ctx.Request.Header.Get("Content-Type")
	boundary := ct[strings.Index(ct, "boundary=")+9:]
	reader := multipart.NewReader(body, boundary)

	xf.debugf("xss filter enter multi-part-form-data\n")

	var bf bytes.Buffer
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		// https://golang.org/src/mime/multipart/multipart_test.go line 230
		bf.WriteString("--")
		bf.WriteString(boundary)
		bf.WriteString("\r\n")

		//val := make([]byte, 0)
		//_, err = part.Read(val)
		var buf bytes.Buffer
		_, err = io.Copy(&buf, part)
		if nil != err {
			return xerror.Wrap(err, "copy body failed")
		}
		val := buf.String()

		bf.WriteString(`Content-Disposition: form-data; name="`)
		bf.WriteString(part.FormName())
		bf.WriteString(`";`)

		if part.FileName() != "" {
			// Content-Disposition: form-data; name="file"; filename="文件.zip"
			bf.WriteString(` filename="`)
			bf.WriteString(part.FileName())
			bf.WriteString("\";\r\n")

			// Content-Type: application/octet-stream
			partCt := part.Header.Get("Content-Type")
			if partCt == "" {
				partCt = `application/octet-stream`
			}
			bf.WriteString("Content-Type: ")
			bf.WriteString(partCt)
			bf.WriteString("\r\n\r\n")
		} else {
			// Content-Disposition: form-data; name="file"
			bf.WriteString("\r\n\r\n")

			if _, ok := xf.skipField[part.FormName()]; !ok {
				val = xf.filterXSS(ctx.FullPath(), part.FormName(), val)
			}
		}

		bf.WriteString(val)
		bf.WriteString("\r\n")
	}

	bf.WriteString("--")
	bf.WriteString(boundary)
	bf.WriteString("--\r\n")

	ctx.Request.Body = ioutil.NopCloser(&bf)
	return nil
}
