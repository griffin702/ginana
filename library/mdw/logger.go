package mdw

import (
	"fmt"
	"ginana/library/log"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"runtime"
	"time"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			// Truncate in a golang < 1.8 safe way
			param.Latency = param.Latency - param.Latency%time.Second
		}
		return fmt.Sprintf("[GiNana] [req] %v |%s %3d %s| %13v | %15s |%s %-7s %s | %s\n%s",
			//时间格式
			param.TimeStamp.Format("2006-01-02 15:04:05"),
			//http请求状态码
			statusColor, param.StatusCode, resetColor,
			//耗时
			param.Latency,
			//客户端IP
			param.ClientIP,
			//http请求方式 get post等
			methodColor, param.Method, resetColor,
			//客户端请求的路径
			param.Path,
			//处理请求错误时设置错误消息
			param.ErrorMessage,
		)
	})
}

// 崩溃恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}
