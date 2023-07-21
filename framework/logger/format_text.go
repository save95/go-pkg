package logger

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type formatText struct {
}

func (lf *formatText) Format(entry *logrus.Entry) ([]byte, error) {

	out := &bytes.Buffer{}
	if entry.Buffer != nil {
		out = entry.Buffer
	}

	out.WriteString(entry.Time.Format("2006-01-02 15:04:05.000"))
	out.WriteString("\t")
	out.WriteString("[")
	out.WriteString(strings.ToUpper(entry.Level.String()))
	out.WriteString("]")
	out.WriteString("\t")
	out.WriteString(entry.Message)

	if len(entry.Data) > 0 {
		out.WriteString("\tDATA:{")
		i := 0
		for key, value := range entry.Data {
			out.WriteString(fmt.Sprintf("%s:%#v", key, value))
			i++
			if i < len(entry.Data) {
				out.WriteString(", ")
			}
		}
		out.WriteString("}")
	}

	out.WriteByte('\n')

	return out.Bytes(), nil
}
