package api

import (
	"github.com/gin-gonic/gin"
)

// RegisterTestRoutes 注册测试路由
func RegisterTestRoutes(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.Engine) {
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/login", UserLogin)
		userGroup.GET("/info", UserInfo)
	}
}

// RegisterGameRoutes 注册游戏相关路由
func RegisterGameRoutes(r *gin.Engine) {
	gameGroup := r.Group("/api/games")
	{
		gameGroup.GET("", GetGameList)
		gameGroup.GET("/:gameCode", GetGameDetail)
	}
}

// RegisterNumberRoutes 注册号码相关路由
func RegisterNumberRoutes(r *gin.Engine) {
	numberGroup := r.Group("/api/numbers")
	{
		numberGroup.POST("/random", GenerateRandomNumbers)
		numberGroup.POST("/save", SaveUserNumber)
		numberGroup.GET("/my", GetMyNumbers)
		numberGroup.PUT("/:id", UpdateUserNumber)
		numberGroup.DELETE("/:id", DeleteUserNumber)
	}
}

// RegisterResultRoutes 注册开奖结果相关路由
func RegisterResultRoutes(r *gin.Engine) {
	resultGroup := r.Group("/api/results")
	{
		resultGroup.GET("/:gameCode", GetDrawResults)
		resultGroup.GET("/:gameCode/:period", GetDrawResultDetail)
	}
}

// RegisterCrawlerRoutes 注册数据抓取相关路由
func RegisterCrawlerRoutes(r *gin.Engine) {
	crawlerGroup := r.Group("/api/crawler")
	{
		crawlerGroup.POST("/crawl/:gameCode", CrawlLatestHandler) // 抓取最新开奖数据
		crawlerGroup.GET("/test/:gameCode", TestCrawlHandler)     // 测试抓取功能
		crawlerGroup.POST("/mock/:gameCode", MockDataHandler)     // 生成模拟数据
	}
}
