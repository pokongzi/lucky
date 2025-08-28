package api

import (
	"net/http"
	"strconv"

	"lucky/common/mysql"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// GetDrawResults 获取开奖结果列表
func GetDrawResults(c *gin.Context) {
	gameCode := c.Param("gameCode")
	if gameCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "游戏代码不能为空",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	results, total, err := service.GetDrawResults(mysql.DB, gameCode, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取开奖结果失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":     results,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// GetDrawResultDetail 获取开奖结果详情
func GetDrawResultDetail(c *gin.Context) {
	gameCode := c.Param("gameCode")
	period := c.Param("period")

	if gameCode == "" || period == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "游戏代码和期号不能为空",
		})
		return
	}

	result, err := service.GetDrawResultByPeriod(mysql.DB, gameCode, period)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "开奖结果不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    result,
	})
}
