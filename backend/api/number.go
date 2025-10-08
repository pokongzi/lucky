package api

import (
	"net/http"
	"strconv"

	"lucky/common/mysql"
	"lucky/model"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// SaveUserNumberRequest 保存用户号码请求
type SaveUserNumberRequest struct {
	GameCode  string            `json:"gameCode" binding:"required"`
	RedBalls  model.NumberArray `json:"redBalls" binding:"required"`
	BlueBalls model.NumberArray `json:"blueBalls" binding:"required"`
	Nickname  string            `json:"nickname"`
	Source    string            `json:"source"`
}

// UpdateUserNumberRequest 更新用户号码请求
type UpdateUserNumberRequest struct {
	Nickname string `json:"nickname"`
	IsActive *bool  `json:"isActive"`
}

// SaveUserNumber 保存用户号码
func SaveUserNumber(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	var req SaveUserNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 获取游戏
	game, err := service.GetGameByCode(mysql.DB, req.GameCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "游戏不存在",
		})
		return
	}

	// 验证号码
	if err := service.ValidateNumbers(game, req.RedBalls, req.BlueBalls); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "号码格式错误",
			"error":   err.Error(),
		})
		return
	}

	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	// 保存用户号码
	userNumber, err := service.SaveUserNumber(mysql.DB, userIDUint, game.ID, req.RedBalls, req.BlueBalls, req.Nickname, req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "保存成功",
		"data":    userNumber,
	})
}

// GetMyNumbers 获取我的号码
func GetMyNumbers(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	gameCode := c.Query("gameCode")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	numbers, total, err := service.GetUserNumbers(mysql.DB, userIDUint, gameCode, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data": gin.H{
			"list":     numbers,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// UpdateUserNumber 更新用户号码
func UpdateUserNumber(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	numberID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID参数错误",
		})
		return
	}

	var req UpdateUserNumberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	err = service.UpdateUserNumber(mysql.DB, userIDUint, numberID, req.Nickname, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

// DeleteUserNumber 删除用户号码
func DeleteUserNumber(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	numberID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "ID参数错误",
		})
		return
	}

	userIDUint, _ := strconv.ParseUint(userID, 10, 64)

	err = service.DeleteUserNumber(mysql.DB, userIDUint, numberID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "删除失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
