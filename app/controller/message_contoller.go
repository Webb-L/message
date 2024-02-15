package controller

import (
	lang "github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"message/app/model"
	"message/app/repository"
	"message/app/request"
	"message/app/response"
	"message/logs"
	"net/http"
)

// MessageIndex 查询消息
//
//	@Summary		查询消息
//	@Description	根据用户凭证查询消息
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Security		BasicAuth
//	@Param			filter		query		string					false	"过滤语句（title = 标题,status = 0|1|2,...）"
//	@Param			sortColumn	query		string					false	"排序列（created_at|updated_at|sender_ids|title|content|category|big_content|introducer_ids|status）"
//	@Param			sortType	query		string					false	"排序类型（asc/desc）"
//	@Param			page		query		int						false	"查询第几页数据"
//	@Success		200			{array}		[]response.Message		"消息信息"
//	@Failure		400			{object}	request.ValidationError	"请求参数错误"
//	@Failure		401			{object}	response.HTTPError		"凭证错误"
//	@Failure		502			{object}	response.HTTPError		"系统异常"
//	@Router			/message [get]
func MessageIndex(ctx *gin.Context) {
	// 从上下文中获取 token
	token, tokenExists := ctx.Get("token")
	// 从上下文中获取 message
	message, messageExists := ctx.Get("message")
	// 从上下文中获取 messageFilters
	messageFilter, messageFilterExists := ctx.Get("messageFilters")

	// 检查 token 和 message 是否存在
	if !tokenExists || !messageExists {
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return
	}

	// 将 token 转换为 MessageToken 类型
	messageToken := token.(model.MessageToken)
	// 将 message 转换为 MessageRequest 类型
	messageRequest := message.(*request.MessageRequest)

	var messageFilters []request.MessageFilterRequest
	if messageFilterExists {
		// 将 messageFilter 转换为 []MessageFilterRequest 类型
		messageFilterRequest := messageFilter.(*[]request.MessageFilterRequest)
		for _, filter := range *messageFilterRequest {
			messageFilters = append(
				messageFilters,
				filter,
			)
		}
	}

	logs.LogInfo.Infof("MessageIndex %v %s", messageRequest, messageToken.AuthId)

	// 返回查询结果
	ctx.JSON(
		http.StatusOK,
		repository.QueryMessagesByMessageTokenMessageRequest(
			messageToken,
			messageRequest,
			messageFilters,
		),
	)
}

// MessageCreate 创建消息
//
//	@Summary		创建消息
//	@Description	创建消息
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Security		BasicAuth
//	@Param			_	body		request.MessageCreateUpdateRequest	true	"创建的数据"
//	@Success		200	{object}	response.Message					"创建成功"
//	@Success		202	{object}	response.HTTPError					"创建失败"
//	@Failure		400	{object}	request.ValidationError				"请求参数错误"
//	@Failure		401	{object}	response.HTTPError					"凭证错误"
//	@Failure		502	{object}	response.HTTPError					"系统异常"
//	@Router			/message [post]
func MessageCreate(ctx *gin.Context) {
	// 从上下文中获取 token
	token, tokenExists := ctx.Get("token")
	// 从上下文中获取 messageCreateUpdate
	messageCreate, messageCreateExists := ctx.Get("messageCreateUpdate")

	// 检查 token 和 messageCreateUpdate 是否存在
	if !tokenExists || !messageCreateExists {
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return
	}

	// 将 token 转换为 MessageToken 类型
	messageToken := token.(model.MessageToken)
	// 将 messageCreateUpdate 转换为 MessageCreateUpdateRequest 类型
	messageCreateRequest := messageCreate.(*request.MessageCreateUpdateRequest)

	// 创建消息
	message, err := repository.CreateMessage(
		messageToken,
		messageCreateRequest,
	)

	if err != nil {
		// 如果创建失败，返回状态码 Accepted
		response.NewError(
			ctx,
			http.StatusAccepted,
			lang.MustGetMessage(ctx, "createMessageFail"),
		)

		logs.LogInfo.Infof("MessageCreate-失败 %s %s", err, messageToken.AuthId)
		return
	}

	logs.LogInfo.Infof("MessageCreate-成功 %s", messageToken.AuthId)

	// 返回创建成功的消息
	ctx.JSON(http.StatusOK, message)
}

// MessageUpdate 更新消息
//
//	@Summary		更新消息
//	@Description	根据消息id更新消息
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Security		BasicAuth
//	@Param			id	path		string								true	"消息id"
//	@Param			_	body		request.MessageCreateUpdateRequest	true	"更新消息"
//	@Success		200	{object}	response.Message					"更新成功"
//	@Success		202	{object}	response.HTTPError					"更新失败"
//	@Failure		400	{object}	request.ValidationError				"请求参数错误"
//	@Failure		401	{object}	response.HTTPError					"凭证错误"
//	@Failure		404	{object}	response.HTTPError					"找不到数据"
//	@Failure		502	{object}	response.HTTPError					"系统异常"
//	@Router			/message/{id} [put]
func MessageUpdate(ctx *gin.Context) {
	// 从上下文中获取 token
	token, tokenExists := ctx.Get("token")
	// 从上下文中获取 messageCreateUpdate
	messageUpdate, messageUpdateExists := ctx.Get("messageCreateUpdate")

	// 检查 token 和 messageCreateUpdate 是否存在
	if !tokenExists || !messageUpdateExists {
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return
	}

	// 将 token 转换为 MessageToken 类型
	messageToken := token.(model.MessageToken)
	// 将 messageCreateUpdate 转换为 MessageCreateUpdateRequest 类型
	messageUpdateRequest := messageUpdate.(*request.MessageCreateUpdateRequest)

	// 根据id查询消息
	oldMessage := repository.QueryMessageById(
		messageToken.AuthId,
		ctx.Param("id"),
	)

	if oldMessage == nil {
		// 如果找不到对应的消息，返回状态码 NotFound
		response.NewError(
			ctx,
			http.StatusNotFound,
			lang.MustGetMessage(ctx, "notFound"),
		)
		return
	}

	// 更新消息
	messageNew, err := repository.UpdateMessage(
		oldMessage,
		messageUpdateRequest,
	)
	if err != nil {
		// 如果更新失败，返回状态码 Accepted
		response.NewError(
			ctx,
			http.StatusAccepted,
			lang.MustGetMessage(ctx, "updateMessageFail"),
		)

		logs.LogInfo.Infof("MessageUpdate-失败 %s %s", err, messageToken.AuthId)
		return
	}

	logs.LogInfo.Infof("MessageUpdate-成功 %s", messageToken.AuthId)

	// 返回更新成功后的消息
	ctx.JSON(http.StatusOK, messageNew)
}

// MessageUpdateStatus 更新消息状态
//
//	@Summary		更新状态
//	@Description	根据数组的数据更新消息状态
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Security		BasicAuth
//	@Param			_	body		[]request.MessageStatusRequest		true	"消息状态"
//	@Success		200	{object}	[]response.MessageStatusResponse	"更新后返回的数据"
//	@Failure		400	{object}	request.ValidationError				"请求参数错误"
//	@Failure		401	{object}	response.HTTPError					"凭证错误"
//	@Failure		502	{object}	response.HTTPError					"系统异常"
//	@Router			/message/status [put]
func MessageUpdateStatus(ctx *gin.Context) {
	// 从上下文中获取 token
	token, tokenExists := ctx.Get("token")
	// 从上下文中获取 messageStatus
	messageStatus, messageExists := ctx.Get("messageStatus")

	// 检查 token 和 messageStatus 是否存在
	if !tokenExists || !messageExists {
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return
	}

	// 将 token 转换为 MessageToken 类型
	messageToken := token.(model.MessageToken)
	// 将 messageStatus 转换为 []MessageStatusRequest 类型
	messageStatusRequest := messageStatus.(*[]request.MessageStatusRequest)

	logs.LogInfo.Infof("MessageUpdateStatus %v %s", messageStatusRequest, messageToken.AuthId)

	// 更新消息状态，并返回更新后的结果
	ctx.JSON(
		http.StatusOK,
		repository.UpdateMessageStatus(
			messageToken,
			messageStatusRequest,
		),
	)
}

// MessageDelete 删除消息
//
//	@Summary		删除消息
//	@Description	根据数组的数据删除消息
//	@Tags			message
//	@Accept			json
//	@Produce		json
//	@Security		BasicAuth
//	@Param			_	body		[]request.MessageDeleteRequest		true	"删除的消息"
//	@Success		200	{object}	[]response.MessageDeleteResponse	"更新后返回的数据"
//	@Failure		400	{object}	request.ValidationError				"请求参数错误"
//	@Failure		401	{object}	response.HTTPError					"凭证错误"
//	@Failure		404	{string}	string								"找不到兑换码"
//	@Failure		502	{string}	string								"系统异常"
//	@Router			/message [delete]
func MessageDelete(ctx *gin.Context) {
	// 从上下文中获取 token
	token, tokenExists := ctx.Get("token")
	// 从上下文中获取 messageDelete
	messageDelete, messageDeleteExists := ctx.Get("messageDelete")

	// 检查 token 和 messageDelete 是否存在
	if !tokenExists || !messageDeleteExists {
		response.NewError(
			ctx,
			http.StatusBadGateway,
			lang.MustGetMessage(ctx, "badGateway"),
		)
		return
	}

	// 将 token 转换为 MessageToken 类型
	messageToken := token.(model.MessageToken)
	// 将 messageDelete 转换为 []MessageDeleteRequest 类型
	messageDeleteRequests := messageDelete.(*[]request.MessageDeleteRequest)

	logs.LogInfo.Infof("MessageDelete %v %s", messageDeleteRequests, messageToken.AuthId)

	// 删除消息，并返回删除结果
	ctx.JSON(
		http.StatusOK,
		repository.DeleteMessagesById(
			messageToken,
			messageDeleteRequests,
		),
	)
}
