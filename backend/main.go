package main

import (
	"log"

	"lucky/api"
	"lucky/common/config"
	"lucky/common/mysql"
	"lucky/migration"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	_ = config.Config

	// 初始化MySQL
	mysql.Init()
	log.Println("MySQL连接成功")

	// 数据库迁移
	if err := migration.AutoMigrate(mysql.DB); err != nil {
		log.Fatal("数据库迁移失败: ", err)
	}

	// 初始化基础数据
	if err := service.InitializeData(mysql.DB); err != nil {
		log.Fatal("基础数据初始化失败: ", err)
	}

	// 初始化Redis
	// redis.Init()

	r := gin.Default()
	gin.SetMode(gin.DebugMode)

	// 注册路由
	api.RegisterTestRoutes(r)
	api.RegisterUserRoutes(r)
	api.RegisterGameRoutes(r)
	api.RegisterNumberRoutes(r)
	api.RegisterResultRoutes(r)
	api.RegisterCrawlerRoutes(r)

	log.Println("服务启动在端口 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("服务启动失败: ", err)
	}
}
