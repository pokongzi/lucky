package api

import (
	"net/http"

	"lucky/common/mysql"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	OpenID    string `json:"openId" binding:"required"`
	Nickname  string `json:"nickname" binding:"required"`
	AvatarURL string `json:"avatarUrl"`
}

// UserLoginResponse 用户登录响应
type UserLoginResponse struct {
	UserID    uint64 `json:"userId"`
	OpenID    string `json:"openId"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
	Token     string `json:"token,omitempty"`
}

// UserLogin 用户登录
func UserLogin(c *gin.Context) {
	var req UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 获取或创建用户
	user, err := service.GetOrCreateUser(mysql.DB, req.OpenID, req.Nickname, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "登录失败",
			"error":   err.Error(),
		})
		return
	}

	// TODO: 生成JWT token
	response := UserLoginResponse{
		UserID:    uint64(user.ID),
		OpenID:    user.OpenID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	// TODO: 从token中获取用户ID
	userID := c.GetHeader("X-User-ID") // 临时使用header传递
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未授权",
		})
		return
	}

	user, err := service.GetUserByID(mysql.DB, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "用户不存在",
		})
		return
	}

	response := UserLoginResponse{
		UserID:    uint64(user.ID),
		OpenID:    user.OpenID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    response,
	})
}
