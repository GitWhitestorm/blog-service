package middleware

import (
	"blogservice/global"
	"blogservice/pkg/logger"
	"bytes"
	"time"

	"github.com/gin-gonic/gin"
)

type AccessLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w AccessLogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}

	return w.ResponseWriter.Write(p)
}

func AccessLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyWrite := &AccessLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyWrite

		beginTime := time.Now().Unix()
		c.Next()
		endtime := time.Now().Unix()

		fields := logger.Fileds{
			"request":  c.Request.PostForm.Encode(),
			"response": bodyWrite.body.String(),
		}

		global.Logger.WithFields(fields).Infof("access log: method:%s,status_code:%d,begin_time:%d,end_time:%d",
			c.Request.Method,
			bodyWrite.Status(),
			beginTime,
			endtime,
		)

	}
}
