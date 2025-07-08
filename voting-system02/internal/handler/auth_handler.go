package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"voting-system/internal/service"
)

type AuthHandler struct {
	service *service.PollService
}

func NewAuthHandler(svc *service.PollService) *AuthHandler {
	return &AuthHandler{service: svc}
}

// loginRequest 用于绑定前端传来的用户名/密码
type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 处理 /api/login 请求
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "用户名或密码不能为空"})
		return
	}

	user, ok := h.service.Authenticate(req.Username, req.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "认证失败"})
		return
	}

	// 把 user.ID 当成“token”返回
	token := user.ID
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"token": token}})
}
