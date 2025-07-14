package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterTestRoutes(r *gin.Engine) {
	r.GET("/api/test", TestHandler)
}

func TestHandler(c *gin.Context) {
	a := 0
	b := 1
	a = a + b
	c.JSON(http.StatusOK, gin.H{"message": a})
}
