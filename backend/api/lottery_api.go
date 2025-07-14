package api

import (
	"lucky/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterLotteryRoutes(r *gin.Engine) {
	r.GET("/api/lottery/random", RandomHandler)
}

func RandomHandler(c *gin.Context) {
	typeStr := c.Query("type")
	var numbers string
	if typeStr == "ssq" {
		numbers = service.RandomSSQ()
	} else if typeStr == "dlt" {
		numbers = service.RandomDLT()
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type参数错误"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"numbers": numbers})
}
