package restful

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/save95/xerror"

	"github.com/gin-gonic/gin"
)

type httpLogger struct {
	ctx        *gin.Context
	respWriter *responseWriter

	retractive string
}

func NewHttpLogger(ctx *gin.Context) *httpLogger {
	hl := &httpLogger{
		ctx:        ctx,
		respWriter: NewResponseWriter(ctx.Writer),
		retractive: "   ",
	}

	ctx.Writer = hl.respWriter

	return hl
}

func (f httpLogger) String() string {
	return fmt.Sprintf(
		"api: %s%s%s%s",
		f.general(),
		f.request(),
		f.response(),
		f.error(),
	)
}

func (f httpLogger) general() string {
	var bf bytes.Buffer
	bf.WriteString("\n[")
	bf.WriteString(f.ctx.Request.Method)
	bf.WriteString("] ")
	bf.WriteString(f.ctx.Request.RequestURI)

	return bf.String()
}

func (f httpLogger) request() string {
	var bs bytes.Buffer
	bs.WriteString("\n\n[Request] ")
	bs.WriteString(f.printHeader(f.ctx.Request.Header))
	bs.WriteString(f.printRequestPayload())

	return bs.String()
}

func (f httpLogger) printHeader(headers http.Header) string {
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

func (f httpLogger) printRequestPayload() string {
	var bf bytes.Buffer
	bf.WriteString("\n [PAYLOAD] ")

	bs, err := ioutil.ReadAll(f.ctx.Request.Body)
	if nil != err {
		bf.WriteByte('\n')
		bf.WriteString(f.retractive)
		bf.WriteString("<read request body: ")
		bf.WriteString(err.Error())
		bf.WriteString(">")

		return bf.String()
	}

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

func (f httpLogger) response() string {
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

func (f httpLogger) error() string {
	errors := f.ctx.Errors.ByType(gin.ErrorTypeAny)
	if len(errors) == 0 {
		return ""
	}

	//err := errors[0].Err
	//if err.IsType(gin.ErrorTypePrivate) {
	//	err = err.Err
	//}

	return f.printError(errors[0].Err)
}

func (f httpLogger) printError(err error) string {
	if nil == err {
		return ""
	}

	var bs bytes.Buffer
	bs.WriteString("\n\n[Error] \n")
	bs.WriteString(f.retractive)
	bs.WriteString(err.Error())
	bs.WriteByte('\n')
	bs.WriteString(" [STACK] \n")

	//stack := fmt.Sprintf("%s%+v", f.retractive, err)
	if xe, ok := err.(xerror.XError); ok {
		err = xe.Unwrap()
	}
	stack := strings.ReplaceAll(fmt.Sprintf("%s%+v", f.retractive, err), "\n", "\n"+f.retractive)
	bs.WriteString(stack)

	return bs.String()
}
