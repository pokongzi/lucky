package api

import (
	"encoding/json"
	"fmt"
	http500 "lucky/common/http/500"
	"lucky/common/redis"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MissingDataRequest 遗漏数据请求参数
type MissingDataRequest struct {
	GameCode    string `json:"gameCode" binding:"required"`    // 游戏代码：dlt、ssq
	PeriodCount int    `json:"periodCount" binding:"required"` // 期数：10、30、50
}

// MissingDataResponse 遗漏数据响应
type MissingDataResponse struct {
	GameCode    string      `json:"gameCode"`
	PeriodCount int         `json:"periodCount"`
	RedBalls    interface{} `json:"redBalls"`
	BlueBalls   interface{} `json:"blueBalls"`
	CachedAt    time.Time   `json:"cachedAt"`
}

// GetMissingData 获取遗漏数据
// @Summary 获取彩票号码遗漏数据
// @Description 支持大乐透(dlt)和双色球(ssq)的遗漏数据查询，支持10期、30期、50期数据
// @Tags 遗漏数据
// @Accept json
// @Produce json
// @Param gameCode query string true "游戏代码" Enums(dlt, ssq)
// @Param periodCount query int true "期数" Enums(10, 30, 50)
// @Success 200 {object} MissingDataResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/missing [get]
func GetMissingData(c *gin.Context) {
	// 获取查询参数
	gameCode := c.Query("gameCode")
	periodCountStr := c.Query("periodCount")

	// 验证参数
	if gameCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gameCode 参数不能为空",
		})
		return
	}

	if gameCode != "dlt" && gameCode != "ssq" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gameCode 只支持 dlt(大乐透) 或 ssq(双色球)",
		})
		return
	}

	periodCount, err := strconv.Atoi(periodCountStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "periodCount 必须是数字",
		})
		return
	}

	if periodCount != 10 && periodCount != 30 && periodCount != 50 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "periodCount 只支持 10、30、50",
		})
		return
	}

	// 构建Redis缓存键
	redisKey := fmt.Sprintf("missing_data:%s:%d", gameCode, periodCount)

	// 尝试从Redis获取缓存数据
	var cachedData MissingDataResponse
	if redis.DB != nil && redis.DB.IsEnabled() {
		err := redis.DB.GetJson(redisKey, &cachedData)
		if err == nil && !cachedData.CachedAt.IsZero() {
			// 检查缓存是否过期（设置为1小时过期）
			if time.Since(cachedData.CachedAt) < time.Hour {
				c.JSON(http.StatusOK, cachedData)
				return
			}
		}
	}

	// 缓存不存在或过期，从500.com获取数据
	var response MissingDataResponse
	response.GameCode = gameCode
	response.PeriodCount = periodCount
	response.CachedAt = time.Now()

	if gameCode == "dlt" {
		// 获取大乐透遗漏数据
		redBalls, blueBalls, err := http500.GetDLTMissingData(periodCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("获取大乐透遗漏数据失败: %v", err),
			})
			return
		}
		response.RedBalls = redBalls
		response.BlueBalls = blueBalls

	} else if gameCode == "ssq" {
		// 获取双色球遗漏数据
		redBalls, blueBalls, err := http500.GetSSQMissingData(periodCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("获取双色球遗漏数据失败: %v", err),
			})
			return
		}
		response.RedBalls = redBalls
		response.BlueBalls = blueBalls
	}

	// 将数据缓存到Redis（设置1小时过期）
	if redis.DB != nil && redis.DB.IsEnabled() {
		responseJSON, _ := json.Marshal(response)
		redis.DB.Set(redisKey, string(responseJSON), time.Hour)
	}

	c.JSON(http.StatusOK, response)
}

// GetMissingDataBatch 批量获取遗漏数据
// @Summary 批量获取多个期数的遗漏数据
// @Description 一次性获取指定游戏的多个期数遗漏数据
// @Tags 遗漏数据
// @Accept json
// @Produce json
// @Param gameCode query string true "游戏代码" Enums(dlt, ssq)
// @Success 200 {object} map[string]MissingDataResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/missing/batch [get]
func GetMissingDataBatch(c *gin.Context) {
	gameCode := c.Query("gameCode")

	// 验证参数
	if gameCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gameCode 参数不能为空",
		})
		return
	}

	if gameCode != "dlt" && gameCode != "ssq" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "gameCode 只支持 dlt(大乐透) 或 ssq(双色球)",
		})
		return
	}

	// 获取10期、30期、50期的数据
	periods := []int{10, 30, 50}
	results := make(map[string]MissingDataResponse)

	for _, period := range periods {
		redisKey := fmt.Sprintf("missing_data:%s:%d", gameCode, period)
		periodKey := fmt.Sprintf("period_%d", period)

		// 尝试从Redis获取缓存数据
		var cachedData MissingDataResponse
		cacheHit := false

		if redis.DB != nil && redis.DB.IsEnabled() {
			err := redis.DB.GetJson(redisKey, &cachedData)
			if err == nil && !cachedData.CachedAt.IsZero() {
				// 检查缓存是否过期（设置为1小时过期）
				if time.Since(cachedData.CachedAt) < time.Hour {
					results[periodKey] = cachedData
					cacheHit = true
				}
			}
		}

		// 如果缓存未命中，从500.com获取数据
		if !cacheHit {
			var response MissingDataResponse
			response.GameCode = gameCode
			response.PeriodCount = period
			response.CachedAt = time.Now()

			if gameCode == "dlt" {
				redBalls, blueBalls, err := http500.GetDLTMissingData(period)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": fmt.Sprintf("获取大乐透%d期遗漏数据失败: %v", period, err),
					})
					return
				}
				response.RedBalls = redBalls
				response.BlueBalls = blueBalls

			} else if gameCode == "ssq" {
				redBalls, blueBalls, err := http500.GetSSQMissingData(period)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": fmt.Sprintf("获取双色球%d期遗漏数据失败: %v", period, err),
					})
					return
				}
				response.RedBalls = redBalls
				response.BlueBalls = blueBalls
			}

			// 缓存数据
			if redis.DB != nil && redis.DB.IsEnabled() {
				responseJSON, _ := json.Marshal(response)
				redis.DB.Set(redisKey, string(responseJSON), time.Hour)
			}

			results[periodKey] = response
		}
	}

	c.JSON(http.StatusOK, results)
}
