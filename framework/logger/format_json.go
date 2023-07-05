package logger

import (
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

type formatJson struct {
}

func (lf *formatJson) Format(entry *logrus.Entry) ([]byte, error) {
	logs := map[string]interface{}{
		"time":    entry.Time.Format("2006-01-02 15:04:05.000"),
		"level":   strings.ToUpper(entry.Level.String()),
		"message": entry.Message,
	}

	if len(entry.Data) > 0 {
		logs["data"] = entry.Data
	}

	return json.Marshal(logs)
}
