package api

import (
	"lucky/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CrawlLatestHandler 抓取最新开奖数据
func CrawlLatestHandler(c *gin.Context) {
	gameCode := c.Param("gameCode")
	if gameCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "游戏代码不能为空",
		})
		return
	}

	crawler := service.NewCrawlerService()
	err := crawler.CrawlAndSaveLatest(gameCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "抓取失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "抓取成功",
	})
}

// TestCrawlHandler 测试抓取功能
func TestCrawlHandler(c *gin.Context) {
	gameCode := c.Param("gameCode")
	if gameCode == "" {
		gameCode = "ssq" // 默认双色球
	}

	crawler := service.NewCrawlerService()
	result, err := crawler.CrawlLatestResults(gameCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "抓取失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "抓取成功",
		"data": result,
	})
}

// MockDataHandler 生成模拟数据
func MockDataHandler(c *gin.Context) {
	gameCode := c.Param("gameCode")
	period := c.Query("period")

	if gameCode == "" {
		gameCode = "ssq"
	}
	if period == "" {
		period = "2025099"
	}

	crawler := service.NewCrawlerService()
	result := crawler.MockDrawResult(gameCode, period)

	err := crawler.SaveDrawResult(result)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "保存失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "模拟数据生成成功",
		"data": result,
	})
}
