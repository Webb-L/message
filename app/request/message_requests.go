package request

import (
	"github.com/gin-gonic/gin"
	"message/logs"
	"strings"
)

type MessageFilterRequest struct {
	Column     string `description:"过滤的列" validate:"required,oneof=created_at updated_at sender_ids title content category big_content introducer_ids status"`
	Comparison string `description:"比较" validate:"required,oneof=> = < >= <= != like in"`
	Value      string `description:"过滤的值" validate:"required"`
}

type MessageRequest struct {
	Filter     string `description:"过滤的语句" form:"filter" example:"status = 1"`
	SortColumn string `description:"排序列" form:"sortColumn" validate:"omitempty,oneof=created_at updated_at sender_ids title content category big_content introducer_ids status" example:"title"`
	SortType   string `description:"排序类型" form:"sortType" validate:"omitempty,oneof=desc asc" example:"asc"`
	Page       int    `description:"查询第几页" form:"page" validate:"omitempty,min=1,max=99999999" example:"1"`
}

// ValidateMessageRequestMiddleware 用于验证消息请求参数的中间件
func ValidateMessageRequestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取 token
		token, _ := ctx.Get("token")
		// 将 token 转换为 MessageToken 类型
		messageToken := token.(string)

		message := &MessageRequest{}
		if !validateStructAndSetContext(
			ctx,
			message,
			"message",
		) {
			logs.LogInfo.Infof("ValidateMessageRequestMiddleware-参数错误 %s", messageToken)
			return
		}

		if message.Filter != "" {
			var filters []MessageFilterRequest
			for _, filterStr := range strings.Split(message.Filter, ",") {
				filter := strings.SplitN(filterStr, " ", 3)
				switch len(filter) {
				case 1:
					filters = append(filters, MessageFilterRequest{
						Column: filter[0],
					})
					break
				case 2:
					filters = append(filters, MessageFilterRequest{
						Column:     filter[0],
						Comparison: filter[1],
					})
					break
				case 3:
					filters = append(filters, MessageFilterRequest{
						Column:     filter[0],
						Comparison: filter[1],
						Value:      filter[2],
					})
					break
				default:
					filters = append(filters, MessageFilterRequest{})
					break
				}
			}

			if err := Validate.Var(&filters, "required,gt=0,dive,required"); err != nil {
				HandlingValidateErrors(ctx, err)
				logs.LogInfo.Infof("ValidateMessageRequestMiddleware-失败-查询过滤语法 %s", messageToken)
				return
			}
			ctx.Set("messageFilters", &filters)
			ctx.Next()
		}

		logs.LogInfo.Infof("ValidateMessageRequestMiddleware-成功 %s", messageToken)
	}
}

type MessageCreateUpdateRequest struct {
	Title         string   `description:"标题" json:"title" validate:"required" example:"标题"`
	Content       string   `description:"简单的内容" json:"content" validate:"required" example:"简单的内容"`
	Category      string   `description:"消息类型" json:"category" validate:"required" example:"important"`
	BigContent    string   `description:"复杂消息" json:"bigContent" validate:"required" example:"复杂的内容"`
	IntroducerIds []string `description:"发给谁" json:"introducerIds" validate:"required,gt=0,dive,required" example:"发给谁"`
}

// ValidateMessageCreateUpdateRequestMiddleware 用于验证创建或更新消息请求参数的中间件
func ValidateMessageCreateUpdateRequestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取 token
		token, _ := ctx.Get("token")
		// 将 token 转换为 MessageToken 类型
		messageToken := token.(string)

		if !validateStructAndSetContext(
			ctx,
			&MessageCreateUpdateRequest{},
			"messageCreateUpdate",
		) {
			logs.LogInfo.Infof("ValidateMessageCreateUpdateRequestMiddleware-失败-参数错误 %s", messageToken)
			return
		}
		logs.LogInfo.Infof("ValidateMessageCreateUpdateRequestMiddleware-成功 %s", messageToken)
	}
}

type MessageStatusRequest struct {
	Id     string `json:"id,omitempty" validate:"required,len=32" example:"1"`
	Status uint8  `json:"status,omitempty" validate:"required,max=2,min=0" example:"1"`
}

// ValidateMessageStatusRequestMiddleware 用于验证消息状态请求参数的中间件
func ValidateMessageStatusRequestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取 token
		token, _ := ctx.Get("token")
		// 将 token 转换为 MessageToken 类型
		messageToken := token.(string)

		if !validateSliceAndSetContext(
			ctx,
			&[]MessageStatusRequest{},
			"messageStatus",
		) {
			logs.LogInfo.Infof("ValidateMessageStatusRequestMiddleware-失败-参数错误 %s", messageToken)
			return
		}
		logs.LogInfo.Infof("ValidateMessageStatusRequestMiddleware-成功 %s", messageToken)
	}
}

type MessageDeleteRequest struct {
	MessageId string `description:"消息id" json:"messageId" validate:"required,len=32" example:"id"`
	Delete    bool   `description:"如果为true表示删除数据否则软删除" json:"delete" example:"false"`
}

// ValidateMessageDeleteRequestMiddleware 用于验证删除消息请求参数的中间件
func ValidateMessageDeleteRequestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取 token
		token, _ := ctx.Get("token")
		// 将 token 转换为 MessageToken 类型
		messageToken := token.(string)

		if !validateSliceAndSetContext(
			ctx,
			&[]MessageDeleteRequest{},
			"messageDelete",
		) {
			logs.LogInfo.Infof("ValidateMessageDeleteRequestMiddleware-失败-参数错误 %s", messageToken)
			return
		}
		logs.LogInfo.Infof("ValidateMessageDeleteRequestMiddleware-成功 %s", messageToken)
	}
}

// ValidateMessageIdRequestMiddleware 用于验证消息ID请求参数的中间件
func ValidateMessageIdRequestMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取 token
		token, _ := ctx.Get("token")
		// 将 token 转换为 MessageToken 类型
		messageToken := token.(string)

		err := Validate.Var(ctx.Param("id"), "required,len=32")
		if err != nil {
			HandlingValidateErrors(ctx, err)
			logs.LogInfo.Infof("ValidateMessageIdRequestMiddleware-失败-参数错误 %s", messageToken)
			return
		}

		ctx.Next()
		logs.LogInfo.Infof("ValidateMessageIdRequestMiddleware-成功 %s", messageToken)
	}
}
