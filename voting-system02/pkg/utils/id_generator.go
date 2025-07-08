package utils

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"time"
)

// GenerateID 生成唯一ID
func GenerateID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		// 如果加密随机数失败，使用纳秒时间戳
		return time.Now().Format("20060102150405")
	}
	return hex.EncodeToString(b)
}

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(c *gin.Context) string {
	userID, exists := c.Get("userID")
	if !exists {
		return "anonymous"
	}
	return userID.(string)
}
