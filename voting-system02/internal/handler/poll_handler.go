package handler

import (
	"net/http"
	_ "time"
	"voting-system/internal/repository"

	"github.com/gin-gonic/gin"

	"voting-system/internal/model"
	"voting-system/internal/service"
	"voting-system/pkg/utils"
)

type PollHandler struct {
	service *service.PollService
}

func NewPollHandler(service *service.PollService) *PollHandler {
	return &PollHandler{service: service}
}

func (h *PollHandler) CreatePoll(c *gin.Context) {
	var req model.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "无效的请求数据",
		})
		return
	}

	poll, err := h.service.CreatePoll(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "创建投票失败",
		})
		return
	}

	c.JSON(http.StatusCreated, model.APIResponse{
		Success: true,
		Data:    poll,
	})
}

func (h *PollHandler) GetAllPolls(c *gin.Context) {
	polls := h.service.GetAllPolls()

	// 简化响应数据
	simplePolls := make([]map[string]interface{}, len(polls))
	for i, poll := range polls {
		simplePolls[i] = map[string]interface{}{
			"id":        poll.ID,
			"title":     poll.Title,
			"createdAt": poll.CreatedAt,
			"endAt":     poll.EndAt,
			"options":   len(poll.Options),
		}
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    simplePolls,
	})
}

func (h *PollHandler) GetPoll(c *gin.Context) {
	id := c.Param("id")

	poll, err := h.service.GetPoll(id)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "投票不存在",
		})
		return
	}

	// 计算投票结果
	voteCounts := make([]int, len(poll.Options))
	for _, optionIndex := range poll.Votes {
		if optionIndex >= 0 && optionIndex < len(poll.Options) {
			voteCounts[optionIndex]++
		}
	}

	// 构建响应
	response := struct {
		Poll       *model.Poll `json:"poll"`
		VoteCounts []int       `json:"vote_counts"`
	}{
		Poll:       poll,
		VoteCounts: voteCounts,
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Data:    response,
	})
}

func (h *PollHandler) Vote(c *gin.Context) {
	pollID := c.Param("id")

	// 简化用户认证
	userID := utils.GetUserIDFromContext(c) // 示例实现

	var req model.VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "无效的请求数据",
		})
		return
	}

	if err := h.service.Vote(pollID, userID, req.OptionID); err != nil {
		switch err.(type) {
		case *repository.RepositoryError:
			c.JSON(http.StatusBadRequest, model.APIResponse{
				Success: false,
				Error:   err.Error(),
			})
		default:
			c.JSON(http.StatusInternalServerError, model.APIResponse{
				Success: false,
				Error:   "投票失败",
			})
		}
		return
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "投票成功",
	})
}
