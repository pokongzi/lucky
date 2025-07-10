package main

import (
	"log"

	"lucky/backend/api"

	"github.com/gin-gonic/gin"

	"lucky/backend/common/config"
	"lucky/backend/common/mysql"
	"lucky/backend/common/redis"
)

func main() {
	// 初始化配置
	_ = config.Config
	// 初始化MySQL
	mysql.Init()
	// 初始化Redis
	redis.Init()

	r := gin.Default()

	api.RegisterUserRoutes(r)
	api.RegisterLotteryRoutes(r)
	api.RegisterTicketRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}
