package main

import (
	"flag"
	"fmt"
	"log"

	"lucky/common/mysql"
	"lucky/service"
)

func main() {
	var (
		gameCode = flag.String("game", "ssq", "游戏代码 (ssq/dlt)")
		action   = flag.String("action", "test", "操作类型 (test/crawl/mock)")
		period   = flag.String("period", "2025099", "期号")
	)
	flag.Parse()

	// 初始化数据库
	mysql.Init()

	crawler := service.NewCrawlerService()

	switch *action {
	case "test":
		fmt.Printf("测试抓取 %s 开奖数据...\n", *gameCode)
		result, err := crawler.CrawlLatestResults(*gameCode)
		if err != nil {
			log.Fatalf("抓取失败: %v", err)
		}
		fmt.Printf("抓取成功: %+v\n", result)

	case "crawl":
		fmt.Printf("抓取并保存 %s 最新开奖数据...\n", *gameCode)
		err := crawler.CrawlAndSaveLatest(*gameCode)
		if err != nil {
			log.Fatalf("抓取保存失败: %v", err)
		}
		fmt.Println("抓取保存成功!")

	case "mock":
		fmt.Printf("生成模拟数据: %s 期号 %s\n", *gameCode, *period)
		result := crawler.MockDrawResult(*gameCode, *period)
		err := crawler.SaveDrawResult(result)
		if err != nil {
			log.Fatalf("保存失败: %v", err)
		}
		fmt.Printf("模拟数据生成成功: %+v\n", result)

	case "schedule":
		fmt.Println("启动定时抓取任务...")
		crawler.ScheduleCrawl()

	default:
		fmt.Printf("不支持的操作: %s\n", *action)
		fmt.Println("支持的操作: test, crawl, mock, schedule")
	}
}
