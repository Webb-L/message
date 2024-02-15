package response

import (
	"message/app/model"
	"time"
)

type Message struct {
	MessageId     string            `json:"message_id" example:"7e55cb38290f49ee2b0e9cfd2adf13e4"`
	SenderIds     model.StringArray `json:"sender_ids" example:"2f14ec370621a8be08c8f0ece459e7e0,22798c5dcd6e5b66c8660c447010d49d,..."`
	Title         string            `json:"title" example:"标题"`
	Content       string            `json:"content" example:"简单的内容"`
	Category      string            `json:"category" example:"important"`
	BigContent    string            `json:"big_content" example:"复杂的内容"`
	IntroducerIds model.StringArray `json:"introducer_ids" example:"fc64c1a807c2e69655f68d31e5caa35d,70c021d35ce60436c115b20b5cf583d0,..."`
	Status        uint8             `json:"status" example:"0"`
	CreatedAt     time.Time         `json:"created_at" example:"2024-02-15T05:49:57Z"`
	UpdatedAt     time.Time         `json:"updated_at" example:"2024-02-15T05:49:57Z"`
}

// MessageStatusResponse 用于封装消息状态更新操作的响应数据
type MessageStatusResponse struct {
	// Id 表示消息的唯一标识符。
	Id string `json:"id"`

	// Status 表示消息的当前状态。
	Status uint8 `json:"status"`

	// Result 表示消息状态更新操作的结果。
	// true 表示更新成功，false 表示更新失败。
	Result bool `json:"result"`
}

// MessageDeleteResponse 表示消息删除操作的响应数据结构
type MessageDeleteResponse struct {
	// Id 表示消息的唯一标识符。
	Id string `json:"id"`

	// Delete 表示消息是否永久删除。
	Delete bool `json:"delete"`

	// Status 表示消息删除操作的状态，用于指示操作是否成功
	Status bool `json:"status"`
}
