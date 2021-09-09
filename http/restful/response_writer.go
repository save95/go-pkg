package restful

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter

	body *bytes.Buffer
}

func NewResponseWriter(w gin.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		body:           &bytes.Buffer{},
	}
}

func (r responseWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
