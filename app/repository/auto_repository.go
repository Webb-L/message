package repository

import (
	"fmt"
	"message/config"
	"message/database"
)

// GetMessageToken 尝试获取一个符合指定条件的MessageToken。
//
// 它返回找到的MessageToken和可能出现的错误。如果记录不存在，将返回nil和gorm.ErrRecordNotFound。
func GetMessageToken(messageToken string) (bool, error) {
	verify := config.AppConfig.App.Verify
	result := database.DB.Table(verify.Table).
		Select(verify.Column).
		Where(fmt.Sprintf("%s = ?", verify.Column), messageToken).
		Limit(1)

	if result.Error != nil {
		// 直接返回错误，包括未找到记录的情况
		return false, result.Error
	}

	// 记录被成功找到，返回token的指针和nil作为错误
	return true, nil
}
