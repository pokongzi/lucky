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
		action   = flag.String("action", "test", "操作类型 (test/crawl/history)")
		pages    = flag.Int("pages", 1, "抓取历史数据的页数")
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

	case "history":
		fmt.Printf("抓取 %s 历史数据，页数：%d...\n", *gameCode, *pages)
		err := crawler.CrawlHistoryByPeriod(*gameCode, *pages)
		if err != nil {
			log.Fatalf("历史数据抓取失败: %v", err)
		}
		fmt.Println("历史数据抓取完成!")

	case "schedule":
		fmt.Println("启动定时抓取任务...")
		crawler.ScheduleCrawl()

	default:
		fmt.Printf("不支持的操作: %s\n", *action)
		fmt.Println("支持的操作: test, crawl, history, schedule")
	}
}
