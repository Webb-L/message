package router

import (
	"github.com/gin-gonic/gin"
	"message/app/controller"
	"message/app/request"
)

// InitMessageRouter 用于初始化消息相关的路由
func InitMessageRouter(router *gin.RouterGroup) {
	// 查询消息
	router.GET(
		"",
		request.ValidateMessageRequestMiddleware(),
		controller.MessageIndex,
	)
	// 新增消息
	router.POST("",
		request.ValidateMessageCreateUpdateRequestMiddleware(),
		controller.MessageCreate,
	)
	// 更新消息
	router.PUT(":id",
		request.ValidateMessageIdRequestMiddleware(),
		request.ValidateMessageCreateUpdateRequestMiddleware(),
		controller.MessageUpdate,
	)
	// 更新消息状态
	router.PUT(
		"status",
		request.ValidateMessageStatusRequestMiddleware(),
		controller.MessageUpdateStatus,
	)
	// 删除通知
	router.DELETE("",
		request.ValidateMessageDeleteRequestMiddleware(),
		controller.MessageDelete,
	)
}
