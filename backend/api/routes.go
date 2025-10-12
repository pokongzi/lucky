package api

import (
	"lucky/middleware"

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
		userGroup.GET("/info", middleware.AuthRequired(), UserInfo)
	}
}

// RegisterAuthRoutes 注册认证相关路由
func RegisterAuthRoutes(r *gin.Engine) {
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/wxlogin", WxLogin)
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
		numberGroup.POST("/save", SaveUserNumber)
		numberGroup.GET("/my", GetMyNumbers)
		numberGroup.PUT("/:id", UpdateUserNumber)
		numberGroup.DELETE("/:id", DeleteUserNumber)
		numberGroup.GET("/:numberId/check", CheckWinning) // 新增：中奖核对
	}
}

// RegisterResultRoutes 注册开奖结果相关路由
func RegisterResultRoutes(r *gin.Engine) {

	// 创建路由组并注册其他开奖结果路由
	resultGroup := r.Group("/api/results")
	{
		resultGroup.GET("/distribution/:gameCode", GetNumberDistribution)
		resultGroup.GET("/:gameCode", GetDrawResults)
		resultGroup.GET("/:gameCode/:period", GetDrawResultDetail) // 通配符路由放在最后
	}
}

// RegisterCrawlerRoutes 注册数据抓取相关路由
func RegisterCrawlerRoutes(r *gin.Engine) {
	crawlerGroup := r.Group("/api/crawler")
	{
		crawlerGroup.POST("/crawl/:gameCode", CrawlLatestHandler) // 抓取最新开奖数据
		crawlerGroup.GET("/test/:gameCode", TestCrawlHandler)     // 测试抓取功能
	}
}

// RegisterMissingRoutes 注册遗漏数据相关路由
func RegisterMissingRoutes(r *gin.Engine) {
	missingGroup := r.Group("/api/missing")
	{
		missingGroup.GET("", GetMissingData)            // 获取指定期数的遗漏数据
		missingGroup.GET("/batch", GetMissingDataBatch) // 批量获取多个期数的遗漏数据
	}
}
