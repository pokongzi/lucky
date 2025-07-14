package main

import (
	"log"

	"lucky/api"

	"github.com/gin-gonic/gin"

	"lucky/common/config"
)

func main() {
	// 初始化配置
	_ = config.Config
	// 初始化MySQL
	// mysql.Init()
	// 初始化Redis
	// redis.Init()

	r := gin.Default()
	gin.SetMode(gin.DebugMode)
	api.RegisterTestRoutes(r)

	// api.RegisterUserRoutes(r)
	// api.RegisterLotteryRoutes(r)
	// api.RegisterTicketRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}
