package xss

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
	"github.com/save95/go-pkg/http/xss"
	"github.com/save95/xerror"
)

type handler struct {
	xssRuleItem

	debug bool

	// 路由特殊规则
	routePolicies map[string]*xssRuleItem
	// 直接跳过的路由
	skipRoutes map[string]struct{}
}

func New(opts ...Option) gin.HandlerFunc {
	xf := &handler{
		xssRuleItem: xssRuleItem{
			skipField: make(map[string]struct{}, 0),
		},
		routePolicies: make(map[string]*xssRuleItem, 0),
		skipRoutes:    make(map[string]struct{}, 0),
	}
	//xf.policy = xf.makePolicy(PolicyStrict)

	for _, opt := range opts {
		opt(xf)
	}

	return xf.filter()
}

func (h *handler) makePolicy(p xss.Policy) *bluemonday.Policy {
	switch p {
	case xss.PolicyNone:
		return nil
	case xss.PolicyStrict:
		return bluemonday.StrictPolicy()
	case xss.PolicyUGC:
		return bluemonday.UGCPolicy()
	default:
		return bluemonday.StrictPolicy()
	}
}

func (h *handler) makeSkipFields(fields []string) map[string]struct{} {
	vals := make(map[string]struct{}, 0)

	for _, field := range fields {
		vals[field] = struct{}{}
	}

	return vals
}

func (h *handler) filter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 指定了忽略路由，直接跳过
		for u := range h.skipRoutes {
			if strings.Contains(ctx.FullPath(), u) {
				h.debugf("xss handler hit skip route, skip\n")
				ctx.Next()
				return
			}
		}

		// 未指定全局规则，而且未指定路由规则，直接跳过
		if h.policy == nil {
			skip := true
			for u := range h.routePolicies {
				if strings.Contains(ctx.FullPath(), u) {
					skip = false
					break
				}
			}
			if skip {
				h.debugf("xss handler no global policy, not hit route rule, skip\n")
				ctx.Next()
				return
			}
		}

		var err error

		switch ctx.Request.Method {
		case http.MethodGet:
			err = h.filterQueryString(ctx)
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := ctx.Request.Header.Get("Content-Type")
			switch ct {
			case "application/json":
				err = h.filterJSON(ctx)
			case "application/x-www-form-urlencoded":
				err = h.filterFormData(ctx)
			default:
				if strings.Contains(ct, "multipart/form-data") {
					err = h.filterMultiPartFormData(ctx)
				}
			}
		}

		if nil != err {
			if xe, ok := err.(xerror.XError); ok {
				err = xe.Unwrap()
			}
			log.Printf("xss handler err: %+v\n", err)
			_ = ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		ctx.Next()
	}
}

func (h *handler) filterXSS(fullPath, key, val string) string {
	for route, item := range h.routePolicies {
		if strings.Contains(fullPath, route) {
			if _, ok := item.skipField[key]; ok {
				h.debugf("xss handler hit route skip field, return origin value\n")
				return val
			}

			fieldPolicy, ok := item.fieldRules[key]
			if ok && fieldPolicy != nil {
				h.debugf("xss handler hit route field rule, return sanitize value\n")
				return fieldPolicy.Sanitize(val)
			}

			if item.policy == nil {
				h.debugf("xss handler hit route rule, none policy, return origin value\n")
				return val
			}

			h.debugf("xss handler hit route rule, return sanitize value\n")
			return item.policy.Sanitize(val)
		}
	}

	if _, ok := h.skipField[key]; ok {
		h.debugf("xss handler hit global skip field, return origin value\n")
		return val
	}

	fieldPolicy, ok := h.fieldRules[key]
	if ok && fieldPolicy != nil {
		h.debugf("xss handler hit global field rule, return sanitize value\n")
		return fieldPolicy.Sanitize(val)
	}

	if h.policy == nil {
		h.debugf("xss handler hit global policy, return sanitize value\n")
		return val
	}

	h.debugf("xss handler hit global policy, return sanitize value\n")
	return h.policy.Sanitize(val)
}

func (h *handler) debugf(format string, vals ...interface{}) {
	if h.debug {
		log.Printf(format, vals...)
	}
}

func (h *handler) filterQueryString(ctx *gin.Context) error {
	params := ctx.Request.URL.Query()
	h.debugf("xss handler input query string: %s\n", params.Encode())
	for key, items := range params {
		params.Del(key)
		for _, val := range items {
			val = h.filterXSS(ctx.FullPath(), key, val)
			if params.Has(key) {
				params.Add(key, val)
			} else {
				params.Set(key, val)
			}
		}
	}

	h.debugf("xss handler output query string: %s\n", params.Encode())
	ctx.Request.URL.RawQuery = params.Encode()
	return nil
}

func (h *handler) filterJSON(ctx *gin.Context) error {
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

	h.debugf("xss handler input json: %s\n", val)

	var data interface{}
	switch val.(type) {
	case map[string]interface{}:
		vals := make(map[string]interface{}, 0)
		for k, v := range val.(map[string]interface{}) {
			vals[k] = h.filterJsonValue(ctx.FullPath(), k, v)
		}
		data = vals
	case []interface{}:
		vals := make([]interface{}, 0)
		for _, v := range val.([]interface{}) {
			vals = append(vals, h.filterJsonValue(ctx.FullPath(), "", v))
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

	h.debugf("xss handler output json: %s\n", bf.String())
	ctx.Request.Body = ioutil.NopCloser(&bf)
	return nil
}

func (h *handler) filterJsonValue(fullPath, key string, val interface{}) interface{} {
	switch val.(type) {
	case map[string]interface{}:
		vals := make(map[string]interface{}, 0)
		for k, v := range val.(map[string]interface{}) {
			vals[k] = h.filterJsonValue(fullPath, key, v)
		}
		return vals
	case []interface{}:
		vals := make([]interface{}, 0)
		for _, v := range val.([]interface{}) {
			vals = append(vals, h.filterJsonValue(fullPath, key, v))
		}
		return vals
	case string:
		return h.filterXSS(fullPath, key, val.(string))
	default:
		return val
	}
}

func (h *handler) filterFormData(ctx *gin.Context) error {
	body := ctx.Request.Body
	if body == nil || body == http.NoBody {
		return nil
	}

	// https://golang.org/src/net/http/httputil/dump.go
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(body); err != nil {
		return xerror.Wrap(err, "read from failed")
	}

	h.debugf("xss handler input x-form-data: %s\n", buf.String())

	m, err := url.ParseQuery(buf.String())
	if err != nil {
		return xerror.Wrap(err, "parse query failed")
	}

	var bf bytes.Buffer
	for key, v := range m {
		val := url.QueryEscape(v[0])
		if _, ok := h.skipField[key]; !ok {
			val = url.QueryEscape(h.filterXSS(ctx.FullPath(), key, v[0]))
		}

		if bf.Len() > 0 {
			bf.WriteByte('&')
		}
		bf.WriteString(key)
		bf.WriteByte('=')
		bf.WriteString(val)
	}

	h.debugf("xss handler output x-form-data: %s\n", bf.String())
	ctx.Request.Body = ioutil.NopCloser(&bf)
	return nil
}

func (h *handler) filterMultiPartFormData(ctx *gin.Context) error {
	body := ctx.Request.Body
	if body == nil || body == http.NoBody {
		return nil
	}

	ct := ctx.Request.Header.Get("Content-Type")
	boundary := ct[strings.Index(ct, "boundary=")+9:]
	reader := multipart.NewReader(body, boundary)

	h.debugf("xss handler enter multi-part-form-data\n")

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

			if _, ok := h.skipField[part.FormName()]; !ok {
				val = h.filterXSS(ctx.FullPath(), part.FormName(), val)
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
