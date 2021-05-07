package restful

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/save95/xerror"
)

type _formatter struct {
	request    *http.Request
	retractive string
	err        error
}

func newRequestError(request *http.Request, err error) *_formatter {
	return &_formatter{
		request:    request,
		retractive: "  ",
		err:        err,
	}
}

func (f _formatter) String() string {
	return fmt.Sprintf(
		"api xerror: \n%s\n%s\n%s\n%s",
		f.uri(),
		f.headers(),
		f.body(),
		f.error(),
	)
}

func (f _formatter) uri() string {
	var bf bytes.Buffer
	bf.WriteString("[")
	bf.WriteString(f.request.Method)
	bf.WriteString("] ")
	bf.WriteString(f.request.RequestURI)

	return bf.String()
}

func (f _formatter) headers() string {
	var bf bytes.Buffer
	bf.WriteString("[HEADER] ")

	for key, val := range f.request.Header {
		bf.WriteByte('\n')
		bf.WriteString(f.retractive)
		bf.WriteString(key)
		bf.WriteString(": ")
		bf.WriteString(strings.Join(val, ", "))
	}

	return bf.String()
}

func (f _formatter) body() string {
	var bf bytes.Buffer
	bf.WriteString("[BODY] ")

	bs, err := ioutil.ReadAll(f.request.Body)
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
	ct := f.request.Header.Get("Content-Type")
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

func (f _formatter) error() string {
	var bs bytes.Buffer
	bs.WriteString("[ERROR] ")
	bs.WriteString(f.err.Error())
	bs.WriteByte('\n')
	bs.WriteString("[STACK] \n")

	stack := fmt.Sprintf("%+v", f.err)
	if xe, ok := f.err.(xerror.XError); ok {
		stack = fmt.Sprintf("%+v", xe.Unwrap())
	}
	bs.WriteString(stack)

	return bs.String()
}
