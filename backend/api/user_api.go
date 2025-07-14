package api

import (
	"lucky/common/mysql"
	"lucky/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthRequest struct {
	OpenID    string `json:"openid" binding:"required"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
}

func RegisterUserRoutes(r *gin.Engine) {
	r.POST("/api/user/auth", AuthHandler)
}

func AuthHandler(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	user, err := service.GetOrCreateUser(mysql.DB, req.OpenID, req.Nickname, req.AvatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户处理失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
