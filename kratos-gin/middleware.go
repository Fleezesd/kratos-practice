package gin

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
)

// GinLogger 基于kratos简易log 所做的middleware
func GinLogger(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		// 记录耗时
		cost := time.Since(start)
		_ = logger.Log(log.LevelInfo,
			"path", path,
			"status", c.Writer.Status(), // rsp 响应状态码
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"ip", c.ClientIP(), // true ip
			"user-agent", c.Request.UserAgent(), // client's user agent
			"errors", c.Errors.ByType(gin.ErrorTypePrivate).String(), // 私有错误
			"cost", cost, // 耗时
		)

	}
}

// GinRecovery
func GinRecovery(logger log.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				// 判断是否为网络断开链接
				// net error
				if ne, ok := err.(*net.OpError); ok {
					// sys error
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
						// 管道中断即连接中断
					}
				}

				// req 复制 会有些信息丢失
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					_ = logger.Log(log.LevelError,
						"path", c.Request.URL.Path,
						"error", err,
						"request", string(httpRequest),
					)

					_ = c.Error(err.(error))
					c.Abort()
					return
				}
				// 其他内部http error 错误
				// 查看是否需要堆栈信息
				if stack {
					_ = logger.Log(log.LevelError,
						"[Recovery from panic]",
						"error", err,
						"request", string(httpRequest),
						"stack", string(debug.Stack()),
					)
				} else {
					_ = logger.Log(log.LevelError,
						"[Recovery from panic]",
						"error", err,
						"request", string(httpRequest),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
