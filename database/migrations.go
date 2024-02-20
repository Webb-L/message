package database

import (
	"fmt"
	"message/app/model"
	"message/config"
	"message/logs"
)

// InitMigration 用于初始化数据库迁移
func InitMigration() {
	// 从配置中获取字符集设置
	charset := config.AppConfig.Database.Params.Character

	// 构建数据库配置字符串
	dbConfig := fmt.Sprintf("charset=%s", charset)

	// 设置表选项为指定的数据库配置，并自动迁移指定的数据模型
	err := DB.Set("gorm:table_options", dbConfig).AutoMigrate(
		// 迁移消息模型
		&model.Message{},
	)
	if err != nil {
		// 输出迁移错误信息
		logs.LogError.Errorf("InitMigration %s", err)
	}
}
