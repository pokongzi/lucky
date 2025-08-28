package api

import (
	"net/http"

	"lucky/common/mysql"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// GetGameList 获取游戏列表
func GetGameList(c *gin.Context) {
	games, err := service.GetActiveGames(mysql.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取游戏列表失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    games,
	})
}

// GetGameDetail 获取游戏详情
func GetGameDetail(c *gin.Context) {
	gameCode := c.Param("gameCode")
	if gameCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "游戏代码不能为空",
		})
		return
	}

	game, err := service.GetGameByCode(mysql.DB, gameCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "游戏不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    game,
	})
}
