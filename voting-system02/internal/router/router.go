package router

import (
	"github.com/gin-gonic/gin"

	"voting-system/internal/config"
	"voting-system/internal/handler"
	"voting-system/internal/repository"
	"voting-system/internal/service"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 初始化依赖
	cfg := config.LoadConfig()
	repo := repository.NewPollRepository(cfg.DataFilePath)
	pollService := service.NewPollService(repo)
	pollHandler := handler.NewPollHandler(pollService)
	authHandler := handler.NewAuthHandler(pollService)

	// 公共路由
	public := r.Group("/api")
	{
		public.POST("/login", authHandler.Login)
		public.GET("/polls", pollHandler.GetAllPolls)
		public.GET("/polls/:id", pollHandler.GetPoll)
	}

	// 需要认证的路由
	protected := r.Group("/api")
	protected.Use(AuthMiddleware(pollService)) // 简化的认证中间件
	{
		protected.POST("/polls", pollHandler.CreatePoll)
		protected.POST("/polls/:id/vote", pollHandler.Vote)
	}

	return r
}

// 简化的认证中间件
func AuthMiddleware(service *service.PollService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "未提供认证令牌"})
			return
		}

		// 简化实现 - 实际应使用JWT验证
		userID := "user123" // 这里应该解析token获取用户ID
		c.Set("userID", userID)
		c.Next()
	}
}
