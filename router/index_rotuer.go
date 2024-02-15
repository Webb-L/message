package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"message/app/middleware"
	"message/config"
	_ "message/docs"
	"net/http"
)

// InitRouter 用于初始化路由配置
func InitRouter(router *gin.Engine) {
	// 添加一个简单的路由示例
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// 创建一个名为 message 的路由组，并应用 AuthMiddleware 中间件
	messageGroup := router.Group("message", middleware.AuthMiddleware())
	InitMessageRouter(messageGroup)

	// 根据配置文件中的设置决定是否允许访问 SwaggerApi
	if config.AppConfig.API.Test {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
