package main

import (
	"log"
	"net/http"

	"lucky/common/mysql"
	"lucky/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// CrawlRequest 抓取请求参数
type CrawlRequest struct {
	GameCode string `json:"game_code" binding:"required"` // 游戏代码 (ssq/dlt)
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// 初始化数据库
	mysql.Init()
	log.Println("MySQL连接成功")

	// 创建 Gin 引擎
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		log.Printf("[%s] %s - 健康检查", c.Request.Method, c.Request.RequestURI)
		c.JSON(http.StatusOK, Response{
			Code:    0,
			Message: "OK",
		})
	})

	// 抓取任务接口
	taskGroup := r.Group("/task")
	{
		// POST /task/crawl - 抓取并保存最新开奖数据
		taskGroup.POST("/crawl", handleCrawl)

		// GET /task/crawl/:gameCode - 抓取指定游戏的最新开奖数据
		taskGroup.GET("/crawl/:gameCode", handleCrawlByGameCode)
	}

	// 启动 gRPC 服务（用于 gocron）
	go startGRPCServer()

	// 启动 HTTP 服务（保留原有抓取接口）
	httpPort := ":8081"
	log.Printf("HTTP 服务启动在端口 %s\n", httpPort)

	if err := r.Run(httpPort); err != nil {
		log.Fatal("HTTP 服务启动失败: ", err)
	}
}

// startGRPCServer 启动 gRPC 服务（支持 HTTP/2）
func startGRPCServer() {
	grpcPort := ":9091"

	// 创建 HTTP 服务器处理 gRPC 请求
	mux := http.NewServeMux()
	mux.HandleFunc("/rpc.Task/Run", handleGRPCRequest)
	mux.HandleFunc("/rpc.Task/Check", handleGRPCRequest)

	// 使用 h2c (HTTP/2 Cleartext) 支持 gRPC over HTTP/2
	h2s := &http2.Server{}
	server := &http.Server{
		Addr:    grpcPort,
		Handler: h2c.NewHandler(mux, h2s),
	}

	log.Printf("gRPC 服务启动在端口 %s (支持 HTTP/2, 用于 gocron 调用)\n", grpcPort)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("gRPC 服务启动失败: %v", err)
	}
}

// handleCrawl 处理 POST 请求的抓取任务
func handleCrawl(c *gin.Context) {
	log.Printf("[%s] %s - 接收到抓取请求", c.Request.Method, c.Request.RequestURI)

	var req CrawlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数错误: %v", err)
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 验证游戏代码
	if req.GameCode != "ssq" && req.GameCode != "dlt" {
		log.Printf("不支持的游戏代码: %s", req.GameCode)
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "不支持的游戏代码，仅支持 ssq 或 dlt",
		})
		return
	}

	crawler := service.NewCrawlerService()

	log.Printf("开始抓取 %s 最新开奖数据...\n", req.GameCode)
	err := crawler.CrawlAndSaveLatest(req.GameCode)
	if err != nil {
		log.Printf("抓取失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "抓取失败: " + err.Error(),
		})
		return
	}

	log.Printf("抓取 %s 成功\n", req.GameCode)
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "抓取成功",
		Data: map[string]string{
			"game_code": req.GameCode,
		},
	})
}

// handleCrawlByGameCode 处理 GET 请求的抓取任务
func handleCrawlByGameCode(c *gin.Context) {
	log.Printf("[%s] %s - 接收到抓取请求", c.Request.Method, c.Request.RequestURI)

	gameCode := c.Param("gameCode")

	// 验证游戏代码
	if gameCode != "ssq" && gameCode != "dlt" {
		log.Printf("不支持的游戏代码: %s", gameCode)
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "不支持的游戏代码，仅支持 ssq 或 dlt",
		})
		return
	}

	crawler := service.NewCrawlerService()

	log.Printf("开始抓取 %s 最新开奖数据...\n", gameCode)
	err := crawler.CrawlAndSaveLatest(gameCode)
	if err != nil {
		log.Printf("抓取失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "抓取失败: " + err.Error(),
		})
		return
	}

	log.Printf("抓取 %s 成功\n", gameCode)
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "抓取成功",
		Data: map[string]string{
			"game_code": gameCode,
		},
	})
}
