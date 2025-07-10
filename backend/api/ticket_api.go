package api

import (
	"lucky/backend/common/mysql"
	"lucky/backend/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SaveTicketRequest struct {
	UserID  uint64 `json:"user_id" binding:"required"`
	Type    int    `json:"type" binding:"required"`
	Numbers string `json:"numbers" binding:"required"`
}

type DeleteTicketRequest struct {
	UserID   uint64 `json:"user_id" binding:"required"`
	TicketID uint64 `json:"ticket_id" binding:"required"`
}

func RegisterTicketRoutes(r *gin.Engine) {
	r.POST("/api/ticket/save", SaveTicketHandler)
	r.GET("/api/ticket/list", ListTicketHandler)
	r.POST("/api/ticket/delete", DeleteTicketHandler)
}

func SaveTicketHandler(c *gin.Context) {
	var req SaveTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	err := service.SaveUserTicket(mysql.DB, req.UserID, req.Type, req.Numbers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "保存成功"})
}

func ListTicketHandler(c *gin.Context) {
	userID, _ := c.GetQuery("user_id")
	typeStr, _ := c.GetQuery("type")
	// 省略参数校验和转换
	// ...
	c.JSON(http.StatusOK, gin.H{"tickets": []interface{}{}}) // TODO: 实现查询逻辑
}

func DeleteTicketHandler(c *gin.Context) {
	var req DeleteTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	err := service.DeleteUserTicket(mysql.DB, req.TicketID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
}
