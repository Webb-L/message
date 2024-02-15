package repository

import (
	"fmt"
	"message/app/model"
	"message/app/request"
	"message/app/response"
	"message/config"
	"message/database"
	"message/utils"
	"slices"
	"strings"
)

// QueryMessagesByMessageTokenMessageRequest 根据消息凭证和消息请求查询消息
func QueryMessagesByMessageTokenMessageRequest(
	// 消息凭证
	token model.MessageToken,
	// 消息请求参数
	messageRequest *request.MessageRequest,
	// 消息过滤器
	filters []request.MessageFilterRequest,
) []response.Message {
	var messages []response.Message
	// 创建消息查询对象
	query := database.DB.Model(&model.Message{})

	// 根据消息凭证中的 AuthId 进行筛选或introducer_ids等于空字符串
	query.Where(
		query.Where(
			"introducer_ids LIKE ?",
			fmt.Sprintf("%%%s%%", token.AuthId),
		).Or("introducer_ids = ?", ""),
	)

	if len(filters) > 0 {
		// 根据传入的过滤器条件进行进一步筛选
		for _, filter := range filters {
			if filter.Comparison == "in" {
				query.Where(
					fmt.Sprintf(
						"%s %s ?",
						filter.Column,
						filter.Comparison,
					),
					strings.Split(filter.Value, "|"),
				)
			} else {
				query.Where(
					fmt.Sprintf(
						"%s %s ?",
						filter.Column,
						filter.Comparison,
					),
					filter.Value,
				)
			}
		}
	}

	// 根据排序字段和排序类型进行排序
	if messageRequest.SortColumn == "" {
		messageRequest.SortColumn = "created_at"
	}
	if messageRequest.SortType == "" {
		messageRequest.SortType = "desc"
	}
	query.Order(fmt.Sprintf(
		"%s %s",
		messageRequest.SortColumn,
		messageRequest.SortType,
	))

	// 根据分页信息查询消息并存储在 messages 中
	maxLimit := config.AppConfig.API.MaxLimit
	if messageRequest.Page == 0 {
		messageRequest.Page = 1
	}
	query.Limit(maxLimit).Offset((messageRequest.Page - 1) * maxLimit).Find(&messages)

	// 返回查询到的消息数组
	return messages
}

// CreateMessage 创建一条新消息
func CreateMessage(
	token model.MessageToken,
	createMessage *request.MessageCreateUpdateRequest,
) (*response.Message, error) {
	messageId := utils.BuildMessageId()
	// 将消息插入到数据库中
	result := database.DB.Model(&model.Message{}).Create(&model.Message{
		// 生成消息 ID
		MessageId: messageId,
		// 设置消息的发送者 ID
		SenderIds: []string{token.AuthId},
		// 设置消息标题
		Title: createMessage.Title,
		// 设置消息内容
		Content: createMessage.Content,
		// 设置消息类别
		Category: createMessage.Category,
		// 设置消息大文本内容
		BigContent: createMessage.BigContent,
		// 设置消息介绍者 ID
		IntroducerIds: createMessage.IntroducerIds,
	})
	// 如果发生错误或者影响的行数为 0，则返回 nil
	if result.Error != nil {
		return nil, result.Error
	}

	newMessage := &response.Message{}
	database.DB.Model(&model.Message{}).Where("message_id = ?", messageId).First(newMessage)
	// 返回创建的消息对象
	return newMessage, nil
}

// UpdateMessage 更新消息内容
func UpdateMessage(
	// 待更新的消息对象
	message *model.Message,
	// 消息更新的内容
	messageUpdate *request.MessageCreateUpdateRequest,
) (*response.Message, error) {
	// 更新消息标题
	message.Title = messageUpdate.Title
	// 更新消息内容
	message.Content = messageUpdate.Content
	// 更新消息类别
	message.Category = messageUpdate.Category
	// 更新消息大文本内容
	message.BigContent = messageUpdate.BigContent
	// 更新消息介绍者 ID
	message.IntroducerIds = messageUpdate.IntroducerIds

	// 保存更新后的消息到数据库中
	result := database.DB.Model(&model.Message{}).Where("id = ?", message.ID).Updates(message)
	// 如果发生错误或者影响的行数为 0，则返回 nil
	if result.Error != nil {
		return nil, result.Error
	}

	newMessage := &response.Message{}
	database.DB.Model(&model.Message{}).Where("id = ?", message.ID).First(newMessage)
	// 返回更新后的消息对象
	return newMessage, nil
}

// UpdateMessageStatus 更新消息状态
func UpdateMessageStatus(
	// 消息凭证
	token model.MessageToken,
	// 要更新的状态请求切片
	status *[]request.MessageStatusRequest,
) []response.MessageStatusResponse {
	// 初始化一个空的MessageStatusResponse切片用于存放每个状态更新的结果
	results := make([]response.MessageStatusResponse, 0)

	// 遍历状态请求切片，对每个请求进行处理
	for _, statusRequest := range *status {
		// 对model.Message模型执行更新操作，设置新的状态
		// 使用LIKE查询匹配introducer_ids，并确保message_id与AuthId相符
		result := database.DB.Model(&model.Message{}).
			Where("introducer_ids LIKE ?", fmt.Sprintf("%%%s%%", token.AuthId)).
			Where("message_id = ?", statusRequest.Id).
			Update("status", statusRequest.Status)

		// 将每次更新操作的结果封装到MessageStatusResponse中，并追加到结果切片中
		results = append(results, response.MessageStatusResponse{
			Id:     statusRequest.Id,
			Status: statusRequest.Status,
			Result: result.Error == nil && result.RowsAffected != 0,
		})
	}

	// 返回所有更新操作的结果
	return results
}

// QueryMessageById 通过消息 ID 查询消息
func QueryMessageById(
	// 用户认证 ID
	authId string,
	// 消息 ID
	id string,
) *model.Message {
	// 初始化一个空的消息对象
	message := &model.Message{}

	// 在数据库中查询匹配条件的消息
	result := database.DB.Model(model.Message{}).
		Where("sender_ids LIKE ?", fmt.Sprintf("%%%s%%", authId)).
		Where("message_id = ?", id).
		First(message)

	// 如果查询出错或者没有匹配到数据，则返回 nil
	if result.Error != nil || result.RowsAffected == 0 {
		return nil
	}

	// 返回查询到的消息对象
	return message
}

// DeleteMessagesById 根据消息 ID 批量删除消息
func DeleteMessagesById(
	// 消息凭证
	token model.MessageToken,
	// 要删除的消息请求切片
	deleteRequests *[]request.MessageDeleteRequest,
) []response.MessageDeleteResponse {
	// 存储要物理删除的消息 ID
	var deletes []string
	// 存储要软删除的消息 ID
	var softDeletes []string

	// 遍历消息删除请求切片，将要删除和要软删除的消息 ID 分别添加到对应的切片中
	for _, messageDelete := range *deleteRequests {
		if messageDelete.Delete {
			deletes = append(deletes, messageDelete.MessageId)
		} else {
			softDeletes = append(softDeletes, messageDelete.MessageId)
		}
	}

	// 物理删除要删除的消息
	database.DB.Unscoped().
		Where("sender_ids LIKE ?", fmt.Sprintf("%%%s%%", token.AuthId)).
		Where("message_id in ?", deletes).
		Delete(&model.Message{})

	// 软删除要软删除的消息
	database.DB.
		Where("sender_ids LIKE ?", fmt.Sprintf("%%%s%%", token.AuthId)).
		Where("message_id in ?", softDeletes).
		Delete(&model.Message{})

	// 存储物理删除失败的消息 ID
	var failDeletes []string
	database.DB.
		Select("message_id").
		Where("sender_ids LIKE ?", fmt.Sprintf("%%%s%%", token.AuthId)).
		Where("message_id in ?", softDeletes).Find(&failDeletes)

	// 存储软删除失败的消息 ID
	var failSoftDeletes []string
	database.DB.
		Select("message_id").
		Where("sender_ids LIKE ?", fmt.Sprintf("%%%s%%", token.AuthId)).
		Where("message_id in ?", softDeletes).Find(&failSoftDeletes)

	// 存储删除操作的结果切片
	var results []response.MessageDeleteResponse

	// 遍历物理删除的消息 ID，将每个 ID 和删除状态封装到 MessageDeleteResponse 中，并根据删除是否成功设置相应的状态
	for _, id := range deletes {
		results = append(results, response.MessageDeleteResponse{
			Id:     id,
			Delete: true,
			Status: slices.Contains(failDeletes, id),
		})
	}

	// 遍历软删除失败的消息 ID，将每个 ID 和删除状态封装到 MessageDeleteResponse 中，并根据删除是否成功设置相应的状态
	for _, id := range failSoftDeletes {
		results = append(results, response.MessageDeleteResponse{
			Id:     id,
			Delete: false,
			Status: slices.Contains(failSoftDeletes, id),
		})
	}

	// 返回删除操作的结果切片
	return results
}
