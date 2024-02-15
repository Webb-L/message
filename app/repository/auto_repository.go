package repository

import (
	"message/app/model"
	"message/database"
)

// GetMessageToken 尝试获取一个符合指定条件的MessageToken。
//
// 它返回找到的MessageToken和可能出现的错误。如果记录不存在，将返回nil和gorm.ErrRecordNotFound。
func GetMessageToken(authId string, authToken string) (model.MessageToken, error) {
	var token model.MessageToken
	result := database.DB.Model(model.MessageToken{}).
		Where("auth_id = ?", authId).
		Where("token = ?", authToken).
		First(&token)

	if result.Error != nil {
		// 直接返回错误，包括未找到记录的情况
		return token, result.Error
	}

	// 记录被成功找到，返回token的指针和nil作为错误
	return token, nil
}
