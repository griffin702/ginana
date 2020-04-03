package mdw

import (
	"ginana/library/log"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"runtime"
	"time"
)

// 日志中间件
func Logger(c *gin.Context) {
	// 开始时间
	start := time.Now()
	// 处理请求
	c.Next()
	// 结束时间
	end := time.Now()
	//执行时间
	latency := end.Sub(start)

	path := c.Request.URL.Path

	clientIP := c.ClientIP()
	method := c.Request.Method
	statusCode := c.Writer.Status()
	log.Infof("%3d|%13v|%15s|%s %s",
		statusCode,
		latency,
		clientIP,
		method, path,
	)
}

// 崩溃恢复中间件
func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			var rawReq []byte
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			if c.Request != nil {
				rawReq, _ = httputil.DumpRequest(c.Request, false)
			}
			log.PrintErrf("%s\n%s\n%v", buf, string(rawReq), err)
			log.Errorf("panic: %v\n", err)
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}

