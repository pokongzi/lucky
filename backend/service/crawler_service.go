package service

import (
"fmt"
"net/http"
"regexp"
"strconv"
"strings"
"time"

"lucky/common/mysql"
"lucky/model"

"github.com/PuerkitoBio/goquery"
"gorm.io/gorm"
)

// DrawDataSource 开奖数据源
type DrawDataSource struct {
Name     string
URL      string
Priority int // 优先级，数字越小优先级越高
}

// CrawlerService 开奖数据抓取服务
type CrawlerService struct {
db      *gorm.DB
sources []DrawDataSource
}

// NewCrawlerService 创建抓取服务实例
func NewCrawlerService() *CrawlerService {
return &CrawlerService{
db: mysql.DB,
sources: []DrawDataSource{
{
Name:     "500彩票网",
URL:      "https://kaijiang.500.com/ssq.shtml",
Priority: 1,
},
{
Name:     "乐彩网",
URL:      "https://www.17500.cn/ssq.html",
Priority: 2,
},
{
Name:     "新浪彩票",
URL:      "https://sports.sina.com.cn/lotto/ssq/",
Priority: 3,
},
},
}
}

// DrawResult 开奖结果数据结构
type DrawResult struct {
Period     string  `json:"period"`      // 期号
DrawDate   string  `json:"draw_date"`   // 开奖日期
RedBalls   []int   `json:"red_balls"`   // 红球号码
BlueBalls  []int   `json:"blue_balls"`  // 蓝球号码
Sales      int64   `json:"sales"`       // 销售额
PoolAmount int64   `json:"pool_amount"` // 奖池金额
GameCode   string  `json:"game_code"`   // 游戏代码
Prizes     []Prize `json:"prizes"`      // 奖项信息
}

// Prize 奖项信息
type Prize struct {
Level       int   `json:"level"`        // 奖级
WinnerNum   int   `json:"winner_num"`   // 中奖注数
WinnerBonus int64 `json:"winner_bonus"` // 单注奖金
}

// CrawlLatestResults 抓取最新开奖结果
func (c *CrawlerService) CrawlLatestResults(gameCode string) (*DrawResult, error) {
for _, source := range c.sources {
result, err := c.crawlFromSource(source, gameCode)
if err != nil {
fmt.Printf("从%s抓取失败: %v\n", source.Name, err)
continue
}
if result != nil {
fmt.Printf("成功从%s获取开奖数据\n", source.Name)
return result, nil
}
}
return nil, fmt.Errorf("所有数据源都抓取失败")
}

// crawlFromSource 从指定数据源抓取
func (c *CrawlerService) crawlFromSource(source DrawDataSource, gameCode string) (*DrawResult, error) {
switch source.Name {
case "500彩票网":
return c.crawlFrom500(gameCode)
case "乐彩网":
return c.crawlFrom17500(gameCode)
case "新浪彩票":
return c.crawlFromSina(gameCode)
default:
return nil, fmt.Errorf("不支持的数据源: %s", source.Name)
}
}

// crawlFrom500 从500彩票网抓取
func (c *CrawlerService) crawlFrom500(gameCode string) (*DrawResult, error) {
if gameCode != "ssq" && gameCode != "dlt" {
return nil, fmt.Errorf("不支持的游戏类型: %s", gameCode)
}

url := fmt.Sprintf("https://kaijiang.500.com/%s.shtml", gameCode)
resp, err := http.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

doc, err := goquery.NewDocumentFromReader(resp.Body)
if err != nil {
return nil, err
}

result := &DrawResult{GameCode: gameCode}

// 解析最新一期数据
doc.Find("table tr").First().Each(func(i int, s *goquery.Selection) {
cells := s.Find("td")
if cells.Length() > 0 {
// 期号
result.Period = strings.TrimSpace(cells.Eq(0).Text())

// 开奖日期
result.DrawDate = strings.TrimSpace(cells.Eq(1).Text())

// 开奖号码
numbersText := strings.TrimSpace(cells.Eq(2).Text())
result.RedBalls, result.BlueBalls = c.parseNumbers(numbersText, gameCode)

// 销售额
salesText := strings.TrimSpace(cells.Eq(3).Text())
result.Sales = c.parseAmount(salesText)

// 奖池金额
poolText := strings.TrimSpace(cells.Eq(4).Text())
result.PoolAmount = c.parseAmount(poolText)
}
})

if result.Period == "" {
return nil, fmt.Errorf("未能解析到期号信息")
}

return result, nil
}

// crawlFrom17500 从乐彩网抓取
func (c *CrawlerService) crawlFrom17500(gameCode string) (*DrawResult, error) {
if gameCode != "ssq" {
return nil, fmt.Errorf("乐彩网暂只支持双色球")
}

url := "https://www.17500.cn/ssq.html"
resp, err := http.Get(url)
if err != nil {
return nil, err
}
defer resp.Body.Close()

doc, err := goquery.NewDocumentFromReader(resp.Body)
if err != nil {
return nil, err
}

result := &DrawResult{GameCode: gameCode}

// 查找最新开奖信息
doc.Find("table tr").Each(func(i int, s *goquery.Selection) {
if i == 1 { // 第二行通常是最新数据
cells := s.Find("td")
if cells.Length() >= 6 {
// 期号
periodLink := cells.Eq(0).Find("a").Text()
result.Period = strings.TrimSpace(periodLink)

// 开奖日期
result.DrawDate = strings.TrimSpace(cells.Eq(1).Text())

// 开奖号码
numbersText := strings.TrimSpace(cells.Eq(2).Text())
result.RedBalls, result.BlueBalls = c.parseNumbers(numbersText, gameCode)

// 销售额
salesText := strings.TrimSpace(cells.Eq(5).Text())
result.Sales = c.parseAmount(salesText)

// 奖池金额
poolText := strings.TrimSpace(cells.Eq(6).Text())
result.PoolAmount = c.parseAmount(poolText)
}
}
})

if result.Period == "" {
return nil, fmt.Errorf("未能解析到期号信息")
}

return result, nil
}

// crawlFromSina 从新浪彩票抓取
func (c *CrawlerService) crawlFromSina(gameCode string) (*DrawResult, error) {
// 新浪彩票的抓取逻辑
// 由于网站结构可能变化，这里提供基础框架
return nil, fmt.Errorf("新浪彩票抓取暂未实现")
}

// parseNumbers 解析号码字符串
func (c *CrawlerService) parseNumbers(numbersText, gameCode string) ([]int, []int) {
var redBalls, blueBalls []int

// 移除多余的空格和特殊字符
numbersText = strings.ReplaceAll(numbersText, " ", "")
numbersText = strings.ReplaceAll(numbersText, "+", " ")

// 使用正则表达式提取数字
re := regexp.MustCompile(`\d+`)
numbers := re.FindAllString(numbersText, -1)

if gameCode == "ssq" { // 双色球
// 前6个是红球，最后1个是蓝球
for i, numStr := range numbers {
if num, err := strconv.Atoi(numStr); err == nil {
if i < 6 {
redBalls = append(redBalls, num)
} else if i == 6 {
blueBalls = append(blueBalls, num)
}
}
}
} else if gameCode == "dlt" { // 大乐透
// 前5个是红球，后2个是蓝球
for i, numStr := range numbers {
if num, err := strconv.Atoi(numStr); err == nil {
if i < 5 {
redBalls = append(redBalls, num)
} else if i < 7 {
blueBalls = append(blueBalls, num)
}
}
}
}

return redBalls, blueBalls
}

// parseAmount 解析金额字符串
func (c *CrawlerService) parseAmount(amountText string) int64 {
// 移除非数字字符
re := regexp.MustCompile(`[^\d]`)
numStr := re.ReplaceAllString(amountText, "")

if amount, err := strconv.ParseInt(numStr, 10, 64); err == nil {
return amount
}
return 0
}

// SaveDrawResult 保存开奖结果到数据库
func (c *CrawlerService) SaveDrawResult(result *DrawResult) error {
// 查找游戏ID
var game model.LotteryGame
err := c.db.Where("name = ?", result.GameCode).First(&game).Error
if err != nil {
return fmt.Errorf("游戏 %s 不存在", result.GameCode)
}

// 检查是否已存在
var existing model.DrawResult
err = c.db.Where("game_id = ? AND draw_number = ?", game.ID, result.Period).First(&existing).Error
if err == nil {
return fmt.Errorf("期号 %s 已存在", result.Period)
}

// 转换日期格式
drawDate, err := time.Parse("2006-01-02", result.DrawDate)
if err != nil {
// 尝试其他日期格式
drawDate, err = time.Parse("2006-1-2", result.DrawDate)
if err != nil {
return fmt.Errorf("日期格式解析失败: %s", result.DrawDate)
}
}

// 转换号码为NumberArray类型
redBalls := model.NumberArray(result.RedBalls)
blueBalls := model.NumberArray(result.BlueBalls)

// 创建数据库记录
drawResult := model.DrawResult{
GameID:     game.ID,
DrawNumber: result.Period,
DrawTime:   drawDate,
RedBalls:   redBalls,
BlueBalls:  blueBalls,
IsActive:   true,
}

return c.db.Create(&drawResult).Error
}

// CrawlAndSaveLatest 抓取并保存最新开奖结果
func (c *CrawlerService) CrawlAndSaveLatest(gameCode string) error {
result, err := c.CrawlLatestResults(gameCode)
if err != nil {
return err
}

return c.SaveDrawResult(result)
}

// CrawlHistoryResults 抓取历史开奖结果
func (c *CrawlerService) CrawlHistoryResults(gameCode string, periods []string) error {
for _, period := range periods {
result, err := c.crawlHistoryByPeriod(gameCode, period)
if err != nil {
fmt.Printf("抓取期号 %s 失败: %v\n", period, err)
continue
}

if err := c.SaveDrawResult(result); err != nil {
fmt.Printf("保存期号 %s 失败: %v\n", period, err)
} else {
fmt.Printf("成功保存期号 %s\n", period)
}

// 控制抓取频率，避免被反爬
time.Sleep(1 * time.Second)
}
return nil
}

// crawlHistoryByPeriod 根据期号抓取历史数据
func (c *CrawlerService) crawlHistoryByPeriod(gameCode, period string) (*DrawResult, error) {
// 这里可以实现具体的历史数据抓取逻辑
// 不同网站的历史数据接口可能不同
return nil, fmt.Errorf("历史数据抓取功能待实现")
}

// ScheduleCrawl 定时抓取任务
func (c *CrawlerService) ScheduleCrawl() {
// 每天定时抓取最新开奖结果
ticker := time.NewTicker(30 * time.Minute) // 每30分钟检查一次
defer ticker.Stop()

for range ticker.C {
fmt.Println("开始定时抓取开奖数据...")

// 抓取双色球
if err := c.CrawlAndSaveLatest("ssq"); err != nil {
fmt.Printf("抓取双色球数据失败: %v\n", err)
}

// 抓取大乐透
if err := c.CrawlAndSaveLatest("dlt"); err != nil {
fmt.Printf("抓取大乐透数据失败: %v\n", err)
}

fmt.Println("定时抓取任务完成")
}
}

// MockDrawResult 生成模拟开奖数据（用于测试）
func (c *CrawlerService) MockDrawResult(gameCode, period string) *DrawResult {
result := &DrawResult{
Period:     period,
DrawDate:   time.Now().Format("2006-01-02"),
GameCode:   gameCode,
Sales:      350000000,
PoolAmount: 2500000000,
}

if gameCode == "ssq" {
// 双色球：6红1蓝
result.RedBalls = []int{5, 8, 13, 17, 18, 29}
result.BlueBalls = []int{2}
} else if gameCode == "dlt" {
// 大乐透：5红2蓝
result.RedBalls = []int{1, 15, 23, 28, 35}
result.BlueBalls = []int{3, 12}
}

return result
}
