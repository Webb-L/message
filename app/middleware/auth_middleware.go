package middleware

import (
	lang "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"message/app/repository"
	"message/app/request"
	"message/app/response"
	"message/logs"
	"net/http"
)

// AuthMiddleware 是一个 Gin 中间件函数，用于验证请求的授权信息。
//
// 该中间件从请求头中获取 Authorization，并解析为授权信息。
//
// 授权信息通过调用 repository.GetMessageToken() 方法获取消息令牌是否有效。
//
// 如果授权信息无效或获取消息令牌失败，则返回相应的错误响应。
//
// 否则，将消息令牌设置到上下文中，并继续处理后续请求。
//
// 返回一个 gin.HandlerFunc 处理程序函数。
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if err := request.Validate.Var(token, "required,len=32"); err != nil {
			logs.LogInfo.Infof("AuthMiddleware-失败 %s %s", err, ctx.ClientIP())
			request.HandlingValidateErrors(ctx, err)
			return
		}

		// 获取消息令牌
		_, err := repository.GetMessageToken(token)
		if err != nil {
			logs.LogInfo.Infof("AuthMiddleware-失败-找不到凭证 %s", ctx.ClientIP())
			response.NewError(
				ctx,
				http.StatusUnauthorized,
				lang.MustGetMessage(ctx, "unauthorized"),
			)
			ctx.Abort()
			return
		}

		logs.LogInfo.Infof("AuthMiddleware-成功 %s %s", ctx.ClientIP(), token)
		// 将消息令牌设置到上下文中
		ctx.Set("token", token)
		ctx.Next()
	}
}
