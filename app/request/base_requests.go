package request

import (
	"errors"
	lang "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"message/app/response"
	"message/logs"
	"net/http"
)

// Validate 是一个用于参数校验的 validator 实例
var Validate = validator.New(validator.WithRequiredStructEnabled())

// ValidationError 用于存储参数校验失败时的相关信息
type ValidationError struct {
	Field   string      `json:"field"`   // 字段名
	Type    string      `json:"type"`    // 字段类型
	Value   interface{} `json:"value"`   // 字段数值
	Param   string      `json:"param"`   // 参数
	Message string      `json:"message"` // 错误消息
}

// validateSliceAndSetContext 对传入的切片类型数据进行校验，并将校验通过的数据存储到 Gin 上下文中
func validateSliceAndSetContext(ctx *gin.Context, object interface{}, saveKey string) bool {
	// 检查并绑定 JSON 数据到对象
	if err := ctx.ShouldBindJSON(object); err != nil {
		logs.LogError.Errorf("validateSliceAndSetContext %s %s %s", saveKey, object, err)
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return false
	}

	// 使用 Validate 对象进行参数校验
	if err := Validate.Var(object, "required,gt=0,dive,required"); err != nil {
		handlingErrors(ctx, err)
		return false
	}

	// 将校验通过的数据存储到 Gin 上下文中
	ctx.Set(saveKey, object)
	ctx.Next()
	return true
}

// validateStructAndSetContext 对传入的结构体数据进行校验，并将校验通过的数据存储到 Gin 上下文中
func validateStructAndSetContext(ctx *gin.Context, object interface{}, saveKey string) bool {
	// 检查并绑定数据到对象
	if err := ctx.ShouldBind(object); err != nil {
		logs.LogError.Errorf("validateStructAndSetContext %s %s %s", saveKey, object, err)
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return false
	}

	// 使用 Validate 对象进行参数校验
	if err := Validate.Struct(object); err != nil {
		handlingErrors(ctx, err)
		return false
	}

	// 将校验通过的数据存储到 Gin 上下文中
	ctx.Set(saveKey, object)
	ctx.Next()
	return true
}

// handlingErrors 处理错误
func handlingErrors(ctx *gin.Context, err error) {
	// 处理校验错误
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		logs.LogError.Errorf("handlingErrors %s", err)
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		ctx.Abort()
		return
	}

	var errorValidations []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		errorValidations = append(errorValidations, ValidationError{
			Field:   err.Field(),
			Type:    err.Type().String(),
			Value:   err.Value(),
			Param:   err.Param(),
			Message: err.Tag(),
		})
	}

	// 返回校验错误信息给客户端
	ctx.JSON(http.StatusBadRequest, errorValidations)
	ctx.Abort()
}
