package response

import "github.com/gin-gonic/gin"

// NewError 是一个示例函数，用于在 Gin 上下文中返回错误响应。
//
// 它接收一个 Gin 上下文对象、状态码和错误信息作为参数。
//
// 创建一个 HTTPError 对象，将状态码和错误信息填充到该对象中。
//
// 最后，使用 JSON 方法将 HTTPError 对象以指定的状态码作为响应返回给客户端。
func NewError(ctx *gin.Context, status int, message string) {
	er := HTTPError{
		Code:    status,
		Message: message,
	}
	ctx.JSON(status, er)
}

// HTTPError 是一个示例结构体，表示 HTTP 错误的信息。
//
// 它有两个字段：Code（状态码）和 Message（错误消息）。
//
// 这些字段将在响应中以 JSON 格式返回给客户端。
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}
