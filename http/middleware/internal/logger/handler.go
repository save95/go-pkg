package logger

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/save95/xerror"

	"github.com/gin-gonic/gin"
)

type handler struct {
	ctx        *gin.Context
	respWriter *responseWriter

	retractive string
}

func New(ctx *gin.Context) ILogger {
	hl := &handler{
		ctx:        ctx,
		respWriter: newResponseWriter(ctx.Writer),
		retractive: "   ",
	}

	ctx.Writer = hl.respWriter

	return hl
}

func (f handler) String() string {
	return fmt.Sprintf(
		"api: %s%s%s%s",
		f.general(),
		f.request(),
		f.response(),
		f.error(),
	)
}

func (f handler) general() string {
	var bf bytes.Buffer
	bf.WriteString("\n[")
	bf.WriteString(f.ctx.Request.Method)
	bf.WriteString("] ")
	bf.WriteString(f.ctx.Request.RequestURI)

	return bf.String()
}

func (f handler) request() string {
	var bs bytes.Buffer
	bs.WriteString("\n\n[Request] ")
	bs.WriteString(f.printHeader(f.ctx.Request.Header))
	bs.WriteString(f.printRequestPayload())

	return bs.String()
}

func (f handler) printHeader(headers http.Header) string {
	var bf bytes.Buffer
	bf.WriteString("\n [HEADER] ")

	for key, val := range headers {
		bf.WriteByte('\n')
		bf.WriteString(f.retractive)
		bf.WriteString(key)
		bf.WriteString(": ")
		bf.WriteString(strings.Join(val, ", "))
	}

	return bf.String()
}

func (f handler) printRequestPayload() string {
	var bf bytes.Buffer

	// 读取 request body 失败，则在日志中显示
	bs, err := ioutil.ReadAll(f.ctx.Request.Body)
	if nil != err {
		bf.WriteString("\n [PAYLOAD] ")
		bf.WriteByte('\n')
		bf.WriteString(f.retractive)
		bf.WriteString("<read body failed: ")
		bf.WriteString(err.Error())
		bf.WriteString(">")

		return bf.String()
	}

	// 如果是 GET 请求，没有 payload 则不显示
	if f.ctx.Request.Method == http.MethodGet && len(bs) == 0 {
		return ""
	}

	bf.WriteString("\n [PAYLOAD] ")

	if len(bs) == 0 {
		bf.WriteByte('\n')
		bf.WriteString(f.retractive)
		bf.WriteString("<nil>")
		return bf.String()
	}

	bf.WriteByte('\n')
	// 通过 header 判断是否为文件上传，
	// 如果是文件，不打印文件内容，仅使用占位符表示
	ct := f.ctx.Request.Header.Get("Content-Type")
	if strings.Contains(ct, "boundary=") {
		boundary := strings.Split(ct, "boundary=")[1]
		reg := regexp.MustCompile(fmt.Sprintf("(%s\r\n.*?filename=[\\s\\S]*?\r\n\r\n)([\\s\\S]*?)(\r\n--%s)", boundary, boundary))
		nstr := reg.ReplaceAllString(string(bs), "$1>>>> FILE DATA <<<<$3")

		bss := strings.Split(nstr, "\r\n")
		for _, s := range bss {
			bf.WriteString(f.retractive)
			bf.WriteString(s)
			bf.WriteString("\r\n")
		}
	} else {
		bf.WriteString(f.retractive)
		bf.Write(bs)
	}

	return bf.String()
}

func (f handler) response() string {
	var bs bytes.Buffer
	bs.WriteString("\n\n[Response] ")
	bs.WriteString("\n [STATUS] ")
	bs.WriteString(strconv.Itoa(f.respWriter.Status()))
	bs.WriteString(f.printHeader(f.respWriter.Header()))

	bs.WriteString("\n [BODY] ")
	bs.WriteByte('\n')
	bs.WriteString(f.retractive)

	body := f.respWriter.body
	if len(body.String()) == 0 {
		bs.WriteString("<nil>")
	} else {
		bs.Write(body.Bytes())
	}

	return bs.String()
}

func (f handler) error() string {
	errs := f.ctx.Errors.ByType(gin.ErrorTypeAny)
	if len(errs) == 0 {
		return ""
	}

	//err := errors[0].Err
	//if err.IsType(gin.ErrorTypePrivate) {
	//	err = err.Err
	//}

	return f.printError(errs[0].Err)
}

func (f handler) printError(err error) string {
	if nil == err {
		return ""
	}

	var bs bytes.Buffer
	bs.WriteString("\n\n[Error] \n")
	bs.WriteString(f.retractive)
	bs.WriteString(err.Error())

	// 如果是 xerror，展示 xfield 内容
	if xf, ok := err.(xerror.XFields); ok {
		fields := xf.GetFields()
		if fields != nil && len(fields) > 0 {
			bs.WriteByte('\n')
			bs.WriteString(" [FIELDS] \n")
			bs.WriteString(f.retractive)

			//jsonIndentStr := f.retractive + f.retractive + f.retractive
			//xfbs, _ := json.MarshalIndent(fields, "", jsonIndentStr)
			xfbs, _ := json.Marshal(fields)
			bs.WriteString(string(xfbs))
		}
	}

	bs.WriteByte('\n')
	bs.WriteString(" [STACK] \n")

	//stack := fmt.Sprintf("%s%+v", f.retractive, err)
	var xe xerror.XError
	if errors.As(err, &xe) {
		err = xe.Unwrap()
	}
	stack := strings.ReplaceAll(fmt.Sprintf("%s%+v", f.retractive, err), "\n", "\n"+f.retractive)
	bs.WriteString(stack)

	return bs.String()
}
