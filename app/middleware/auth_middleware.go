package middleware

import (
	"encoding/base64"
	"fmt"
	lang "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"message/app/repository"
	"message/app/response"
	"message/logs"
	"net/http"
	"strings"
)

// parseAuthorization 函数用于解析 Authorization 头中的 Basic token。
//
// 它接受一个字符串参数 basic，该字符串应包含完整的 Authorization 请求头的值。
//
// 函数返回一个包含两个字符串的切片和一个错误对象。
//
// 如果解析成功，切片中将包含从基于 Base64 编码的 token 中解码得到的信息；
//
// 如果解析失败，则返回一个错误。
func parseAuthorization(authorization string) (token []string, err error) {
	// 使用空格将传入的 basic 字符串分割成两部分。
	parts := strings.Split(authorization, " ")
	// 检查分割后的结果是否正好两部分，并且第一部分（不区分大小写）是否为"basic"。
	if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
		// 如果不满足条件，返回一个空的字符串切片和一个格式错误的错误信息。
		return []string{}, fmt.Errorf("格式错误")
	}

	// 尝试对第二部分（即 Base64 编码的 token）进行解码。
	decodeByte, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		// 如果解码过程中发生错误，同样返回一个空的字符串切片和一个格式错误的错误信息。
		return []string{}, fmt.Errorf("格式错误")
	}

	// 将解码后的字节序列转换为字符串，并以冒号为分隔符进行分割。
	info := strings.Split(string(decodeByte), ":")
	// 检查分割后的结果是否正好两部分，这通常对应于用户名和密码。
	if len(info) != 2 {
		// 如果不满足条件，再次返回一个空的字符串切片和一个格式错误的错误信息。
		return []string{}, fmt.Errorf("格式错误")
	}
	// 如果所有检查都通过，则返回解析得到的信息和 nil 错误。
	return info, nil
}

// AuthMiddleware 是一个 Gin 中间件函数，用于验证请求的授权信息。
//
// 该中间件从请求头中获取 Authorization，并解析为授权信息。
//
// 授权信息通过调用 repository.GetMessageToken() 方法获取消息令牌。
//
// 如果授权信息无效或获取消息令牌失败，则返回相应的错误响应。
//
// 否则，将消息令牌设置到上下文中，并继续处理后续请求。
//
// 返回一个 gin.HandlerFunc 处理程序函数。
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头中获取 Authorization 并解析 Authorization
		authorization, err := parseAuthorization(ctx.GetHeader("Authorization"))
		if err != nil {
			logs.LogInfo.Infof("AuthMiddleware-失败 %s %s", err, ctx.ClientIP())
			response.NewError(
				ctx,
				http.StatusUnauthorized,
				lang.MustGetMessage(ctx, "unauthorized"),
			)
			ctx.Abort()
			return
		}

		// 获取消息令牌
		token, err := repository.GetMessageToken(authorization[0], authorization[1])
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

		logs.LogInfo.Infof("AuthMiddleware-成功 %s %s", ctx.ClientIP(), token.AuthId)
		// 将消息令牌设置到上下文中
		ctx.Set("token", token)
		ctx.Next()
	}
}
