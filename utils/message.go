package utils

import (
	"crypto/md5"
	"fmt"
	"github.com/google/uuid"
)

// BuildMessageId 用于生成消息ID
func BuildMessageId() string {
	// 生成一个新的UUID
	u := uuid.New()
	// 将UUID转换为字节切片
	uuidBytes := []byte(u.String())
	// 使用MD5算法加密UUID
	md5Bytes := md5.Sum(uuidBytes)
	// 将加密后的字节切片转换为16进制字符串
	md5String := fmt.Sprintf("%x", md5Bytes)

	return md5String
}
