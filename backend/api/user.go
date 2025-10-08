package api

import (
	"net/http"

	"lucky/common/jwt"
	"lucky/common/mysql"
	"lucky/middleware"
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

	// 生成JWT token
	token, err := jwt.GenerateToken(uint64(user.ID), user.OpenID, user.Nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成token失败",
			"error":   err.Error(),
		})
		return
	}

	response := UserLoginResponse{
		UserID:    uint64(user.ID),
		OpenID:    user.OpenID,
		Nickname:  user.Nickname,
		AvatarURL: user.AvatarURL,
		Token:     token,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}

// UserInfo 获取用户信息
func UserInfo(c *gin.Context) {
	// 从中间件注入的上下文获取用户
	if userModel, ok := middleware.GetCurrentUser(c); ok {
		response := UserLoginResponse{
			UserID:    uint64(userModel.ID),
			OpenID:    userModel.OpenID,
			Nickname:  userModel.Nickname,
			AvatarURL: userModel.AvatarURL,
		}

		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "success",
			"data":    response,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"code":    401,
		"message": "未授权",
	})
}
