package utils

import (
	"github.com/gin-gonic/gin"
	"message/logs"
	"time"
)

// AccessLogger 是自定义的日志记录器中间件
func AccessLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 记录请求开始时间
		startTime := time.Now()

		// 调用后续的中间件和处理函数
		ctx.Next()

		// 记录请求结束时间
		endTime := time.Now()
		// 计算请求耗时
		latency := endTime.Sub(startTime)

		// 将日志信息输出到标准输出流
		logs.LogAccess.Infof(
			"%s    %s    %s    %s",
			ctx.Request.Method,
			ctx.Request.URL.Path,
			latency.String(),
			ctx.ClientIP(),
		)
	}
}
