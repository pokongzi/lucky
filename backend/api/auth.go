package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"lucky/common/config"
	"lucky/common/jwt"
	"lucky/common/mysql"
	"lucky/model"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// WxLoginRequest 微信小程序登录请求
type WxLoginRequest struct {
	Code      string `json:"code" binding:"required"`
	Nickname  string `json:"nickname"`  // 用户昵称
	AvatarURL string `json:"avatarUrl"` // 用户头像
}

// WxLoginResponse 登录响应
type WxLoginResponse struct {
	Token     string     `json:"token"`
	ExpiresAt time.Time  `json:"expiresAt"`
	User      WxUserInfo `json:"user"`
}

// WxUserInfo 微信用户信息
type WxUserInfo struct {
	ID        int64  `json:"id"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatarUrl"`
}

type wxCode2SessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// WxLogin 微信小程序登录
func WxLogin(c *gin.Context) {
	var req WxLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 获取微信配置
	section := config.Config.Section("wechat")
	appid := section.Key("app_id").String()
	secret := section.Key("app_secret").String()
	if appid == "" || secret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "未配置微信小程序appid/secret",
		})
		return
	}

	// 调用微信接口获取openid
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", appid, secret, req.Code)
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    502,
			"message": "请求微信接口失败",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var wx wxCode2SessionResp
	if err := json.NewDecoder(resp.Body).Decode(&wx); err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code":    502,
			"message": "解析微信响应失败",
			"error":   err.Error(),
		})
		return
	}

	if wx.ErrCode != 0 || wx.OpenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "微信登录失败",
			"error":   wx.ErrMsg,
		})
		return
	}

	// 获取数据库连接
	db := mysql.DB

	// 设置默认用户信息
	nickname := req.Nickname
	avatarURL := req.AvatarURL
	if nickname == "" {
		nickname = "微信用户"
	}

	// 获取或创建用户
	user, err := service.GetOrCreateUser(db, wx.OpenID, nickname, avatarURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "用户信息处理失败",
			"error":   err.Error(),
		})
		return
	}

	// 更新登录信息
	userDAO := model.NewUserDAO(db)
	clientIP := c.ClientIP()
	now := time.Now()
	if err := userDAO.UpdateLoginInfo(user.ID, now, clientIP); err != nil {
		// 记录日志但不影响登录流程
		fmt.Printf("更新登录信息失败: %v\n", err)
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

	// 计算token过期时间
	expiresAt := time.Now().Add(jwt.GetAccessTokenExpire())

	// 返回登录成功响应
	response := WxLoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User: WxUserInfo{
			ID:        user.ID,
			Nickname:  user.Nickname,
			AvatarURL: user.AvatarURL,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data":    response,
	})
}
