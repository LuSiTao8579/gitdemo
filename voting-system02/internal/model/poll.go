package model

import (
	"time"
)

// 投票模型
type Poll struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Options     []PollOption   `json:"options"`
	CreatedAt   time.Time      `json:"created_at"`
	EndAt       time.Time      `json:"end_at"`
	Votes       map[string]int `json:"votes"` // userID -> optionIndex
}

// 投票选项
type PollOption struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// 用户模型
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// API响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 创建投票请求
type CreatePollRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Options     []string  `json:"options" binding:"required,min=2"`
	EndAt       time.Time `json:"end_at" binding:"required"`
}

// 投票请求
type VoteRequest struct {
	OptionID string `json:"option_id" binding:"required"`
}
