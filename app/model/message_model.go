package model

import (
	"gorm.io/gorm"
)

// 定义消息状态的常量
const (
	Unread   = iota // 未读状态，默认为 0
	Read            // 已读状态
	Archived        // 归档状态
)

// Message 消息
type Message struct {
	gorm.Model    `json:"-"`
	MessageId     string      `gorm:"type:varchar(32);index;unique;not null;comment:消息id"`
	SenderIds     StringArray `gorm:"type:text;comment:发送者的ID集合"`
	Title         string      `gorm:"type:varchar(25);not null;comment:消息标题"`
	Content       string      `gorm:"type:varchar(50);not null;comment:消息内容"`
	Category      string      `gorm:"type:varchar(50);index;not null;comment:消息类别"`
	BigContent    string      `gorm:"type:longtext;not null;comment:消息的详细内容"`
	IntroducerIds StringArray `gorm:"type:text;comment:接收者的ID集合"`
	Status        uint8       `gorm:"type:tinyint;default:0;comment:消息阅读状态"`
}

type MessageCategory struct {
}
