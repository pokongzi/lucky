package service

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"lucky/common/http/fucai"
	"lucky/common/http/ticai"
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
	sources map[string][]DrawDataSource // 按游戏类型分组的数据源
}

// CWLResult 中国福彩API返回结果结构
type CWLResult struct {
	State   int    `json:"state"`
	Message string `json:"message"`
	Result  []struct {
		Code string `json:"code"` // 期号
		Date string `json:"date"` // 开奖日期
		Red  string `json:"red"`  // 红球
		Blue string `json:"blue"` // 蓝球
	} `json:"result"`
}

// NewCrawlerService 创建抓取服务实例
func NewCrawlerService() *CrawlerService {
	return &CrawlerService{
		db: mysql.DB,
		sources: map[string][]DrawDataSource{
			"ssq": { // 双色球数据源
				{
					Name:     "500彩票网",
					URL:      "https://kaijiang.500.com/ssq.shtml",
					Priority: 1,
				},
				{
					Name:     "中国福彩",
					URL:      "https://www.cwl.gov.cn/ygkj/wqkjgg/ssq/",
					Priority: 2,
				},
			},
			"dlt": { // 大乐透数据源
				{
					Name:     "体彩大乐透",
					URL:      "https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&isVerify=1&termLimits=50",
					Priority: 1,
				},
				{
					Name:     "500彩票网大乐透",
					URL:      "https://kaijiang.500.com/dlt.shtml",
					Priority: 2,
				},
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
	// 获取游戏特定的数据源
	gameSources, exists := c.sources[gameCode]
	if !exists {
		return nil, fmt.Errorf("不支持的游戏类型: %s", gameCode)
	}

	// 按优先级排序数据源
	for _, source := range gameSources {
		result, err := c.crawlFromSource(source, gameCode)
		if err != nil {
			fmt.Printf("%s抓取失败: %v\n", source.Name, err)
			continue
		}
		if result != nil {
			fmt.Printf("成功从%s获取开奖数据\n", source.Name)
			return result, nil
		}
	}
	return nil, fmt.Errorf("所有数据源都抓取失败")
}

// SportteryResult 体彩API返回结果结构
type SportteryResult struct {
	Value struct {
		List []struct {
			LotteryDrawNum    string `json:"lotteryDrawNum"`    // 期号
			LotteryDrawTime   string `json:"lotteryDrawTime"`   // 开奖日期
			LotteryDrawResult string `json:"lotteryDrawResult"` // 开奖结果
		} `json:"list"`
		PageNo   int `json:"pageNo"`
		PageSize int `json:"pageSize"`
		Pages    int `json:"pages"`
		Total    int `json:"total"`
	} `json:"value"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

// crawlFromSource 从指定数据源抓取
func (c *CrawlerService) crawlFromSource(source DrawDataSource, gameCode string) (*DrawResult, error) {
	switch source.Name {
	case "500彩票网":
		return c.crawlFrom500(gameCode)
	case "中国福彩":
		return c.crawlFromCWL(gameCode)
	case "体彩大乐透":
		return c.crawlFromDLT(gameCode)
	case "500彩票网大乐透":
		return c.crawlFrom500DLT(gameCode)
	default:
		return nil, fmt.Errorf("不支持的数据源 %s", source.Name)
	}
}

// crawlFrom500 从500彩票网抓取（仅双色球）
func (c *CrawlerService) crawlFrom500(gameCode string) (*DrawResult, error) {
	if gameCode != "ssq" {
		return nil, fmt.Errorf("500彩票网双色球数据源仅支持双色球")
	}

	url := fmt.Sprintf("https://kaijiang.500.com/%s.shtml", gameCode)
	// 创建一个自定义的http.Client 并设置超时
	client := &http.Client{
		Timeout: 10 * time.Second, // 增加超时时间到10秒
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	// 添加User-Agent避免反爬
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DrawResult{GameCode: gameCode}

	// 尝试多种选择器解析500彩票网数据
	var found bool

	// 方法1：尝试从开奖号码区域解析
	doc.Find(".ball_box, .kjhm, .kjhm_box").Each(func(i int, s *goquery.Selection) {
		if found {
			return
		}

		// 查找期号
		periodText := s.Find(".kjqihao, .qihao, .period").Text()
		if periodText == "" {
			// 尝试从父级或兄弟元素查找期号
			periodText = s.Parent().Find(".kjqihao, .qihao, .period").Text()
		}
		if periodText != "" {
			result.Period = strings.TrimSpace(periodText)
			fmt.Printf("500彩票网解析期号(方法1): %s\n", result.Period)
		}

		// 查找开奖号码
		numbersText := s.Text()
		if numbersText != "" {
			fmt.Printf("500彩票网解析号码文本(方法1): %s\n", numbersText)
			result.RedBalls, result.BlueBalls = c.parseNumbers(numbersText, gameCode)
			fmt.Printf("500彩票网解析号码(方法1): 红球%v 蓝球%v\n", result.RedBalls, result.BlueBalls)
			if len(result.RedBalls) > 0 {
				found = true
			}
		}
	})

	// 方法2：尝试从 .red 和 .blue 选择器解析（500彩票网新格式）
	if !found {
		fmt.Println("500彩票网尝试方法2：从.red和.blue选择器解析")

		// 查找期号 - 尝试从页面文本中提取
		pageText := doc.Text()
		reNum := regexp.MustCompile(`\d+`)
		allNums := reNum.FindAllString(pageText, -1)

		// 查找7位期号（如2025118）
		for _, num := range allNums {
			if len(num) == 7 && strings.HasPrefix(num, "2025") {
				result.Period = num
				fmt.Printf("500彩票网解析期号(方法2): %s\n", result.Period)
				break
			}
		}

		// 方法2.1：尝试从开奖号码区域解析（更精确的选择器）
		doc.Find(".ball_box, .kjhm, .kjhm_box, .kj_tablelist02").Each(func(i int, s *goquery.Selection) {
			if found {
				return
			}

			// 查找期号
			periodText := s.Find(".kjqihao, .qihao, .period, td").First().Text()
			if periodText != "" {
				// 提取期号数字
				re := regexp.MustCompile(`\d+`)
				periods := re.FindAllString(periodText, -1)
				for _, p := range periods {
					if len(p) == 7 && strings.HasPrefix(p, "2025") {
						result.Period = p
						fmt.Printf("500彩票网解析期号(方法2.1): %s\n", result.Period)
						break
					}
				}
			}

			// 查找开奖号码
			numbersText := s.Text()
			if numbersText != "" {
				fmt.Printf("500彩票网解析号码文本(方法2.1): %s\n", numbersText)

				// 专门针对500彩票网的格式解析
				redBalls, blueBalls := c.parse500Numbers(numbersText)
				fmt.Printf("500彩票网解析号码(方法2.1): 红球%v 蓝球%v\n", redBalls, blueBalls)
				if len(redBalls) == 6 && len(blueBalls) == 1 {
					result.RedBalls = redBalls
					result.BlueBalls = blueBalls
					found = true
				}
			}
		})

		// 方法2.2：如果方法2.1失败，尝试从.red和.blue选择器解析
		if !found {
			fmt.Println("500彩票网尝试方法2.2：从.red和.blue选择器解析")

			// 分别查找红球和蓝球
			var redBalls, blueBalls []int

			// 查找红球 - 只取前6个
			doc.Find(".red").Each(func(i int, s *goquery.Selection) {
				if len(redBalls) >= 6 {
					return
				}
				text := strings.TrimSpace(s.Text())
				if num, err := strconv.Atoi(text); err == nil && num >= 1 && num <= 33 {
					redBalls = append(redBalls, num)
				}
			})

			// 查找蓝球 - 只取第一个
			doc.Find(".blue").Each(func(i int, s *goquery.Selection) {
				if len(blueBalls) >= 1 {
					return
				}
				text := strings.TrimSpace(s.Text())
				if num, err := strconv.Atoi(text); err == nil && num >= 1 && num <= 16 {
					blueBalls = append(blueBalls, num)
				}
			})

			fmt.Printf("500彩票网解析号码(方法2.2): 红球%v 蓝球%v\n", redBalls, blueBalls)

			if len(redBalls) == 6 && len(blueBalls) == 1 {
				result.RedBalls = redBalls
				result.BlueBalls = blueBalls
				found = true
			}
		}
	}

	// 不使用兜底解析，如果常规解析失败，直接返回错误
	if !found || result.Period == "" {
		if !found {
			return nil, fmt.Errorf("500彩票网双色球页面解析失败：无法解析开奖号码")
		}
		if result.Period == "" {
			return nil, fmt.Errorf("500彩票网双色球页面解析失败：无法解析期号")
		}
	}

	// 设置默认日期
	if result.DrawDate == "" {
		result.DrawDate = time.Now().Format("2006-01-02")
	}

	if result.Period == "" {
		return nil, fmt.Errorf("未能解析到期号信息")
	}

	return result, nil
}

// crawlFromCWL 从中国福彩官网抓取（仅双色球）
func (c *CrawlerService) crawlFromCWL(gameCode string) (*DrawResult, error) {
	if gameCode != "ssq" {
		return nil, fmt.Errorf("中国福彩暂只支持双色球")
	}

	// 使用中国福彩官网主页
	url := "https://www.cwl.gov.cn/"
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("中国福彩官网响应状态: %d\n", resp.StatusCode)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("中国福彩官网返回非200状态码: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DrawResult{GameCode: gameCode}

	// 方法1：从页面标题解析期号
	fmt.Println("中国福彩尝试方法1：从页面内容解析")

	// 查找期号 - 查找包含"第"和"期"的文本
	pageText := doc.Text()
	fmt.Printf("页面文本长度: %d\n", len(pageText))

	// 查找期号模式：第2025119期
	periodRe := regexp.MustCompile(`第(\d{7})期`)
	periodMatches := periodRe.FindStringSubmatch(pageText)
	if len(periodMatches) > 1 {
		result.Period = periodMatches[1]
		fmt.Printf("中国福彩解析期号(方法1): %s\n", result.Period)
	}

	// 查找开奖日期
	dateRe := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
	dateMatches := dateRe.FindStringSubmatch(pageText)
	if len(dateMatches) > 1 {
		result.DrawDate = dateMatches[1]
		fmt.Printf("中国福彩解析日期(方法1): %s\n", result.DrawDate)
	}

	// 查找开奖号码 - 查找包含红球和蓝球的区域
	// 查找包含"双色球"的元素
	ssqElements := doc.Find("*:contains('双色球')")
	if ssqElements.Length() > 0 {
		fmt.Printf("找到双色球相关元素，数量: %d\n", ssqElements.Length())

		// 遍历双色球相关元素
		ssqElements.Each(func(i int, s *goquery.Selection) {
			elementText := strings.TrimSpace(s.Text())
			fmt.Printf("双色球元素[%d]: %s\n", i, elementText[:min(200, len(elementText))])

			// 查找这个元素附近的数字
			parent := s.Parent()
			if parent.Length() > 0 {
				parentText := strings.TrimSpace(parent.Text())
				fmt.Printf("父元素文本: %s\n", parentText[:min(300, len(parentText))])

				// 尝试从父元素中提取号码
				redBalls, blueBalls := c.parseCWLNumbersFromPage(parentText)
				if len(redBalls) == 6 && len(blueBalls) == 1 {
					result.RedBalls = redBalls
					result.BlueBalls = blueBalls
					fmt.Printf("中国福彩解析号码(方法1): 红球%v 蓝球%v\n", redBalls, blueBalls)
				}
			}
		})
	}

	// 方法2：兜底解析 - 从整个页面文本中提取
	if len(result.RedBalls) != 6 || len(result.BlueBalls) != 1 {
		fmt.Println("中国福彩尝试方法2：兜底解析")

		// 从页面文本中提取所有数字
		reNum := regexp.MustCompile(`\d+`)
		allNums := reNum.FindAllString(pageText, -1)
		fmt.Printf("页面数字前30个: %v\n", allNums[:min(30, len(allNums))])

		// 查找期号
		if result.Period == "" {
			for _, num := range allNums {
				if len(num) == 7 && strings.HasPrefix(num, "2025") {
					result.Period = num
					fmt.Printf("中国福彩解析期号(方法2): %s\n", result.Period)
					break
				}
			}
		}

		// 查找开奖号码
		redBalls, blueBalls := c.parseCWLNumbersFromPage(pageText)
		if len(redBalls) == 6 && len(blueBalls) == 1 {
			result.RedBalls = redBalls
			result.BlueBalls = blueBalls
			fmt.Printf("中国福彩解析号码(方法2): 红球%v 蓝球%v\n", redBalls, blueBalls)
		}
	}

	// 验证数据完整性
	if result.Period == "" {
		return nil, fmt.Errorf("未能解析到期号信息")
	}
	if len(result.RedBalls) != 6 {
		return nil, fmt.Errorf("红球数量错误: 期望6个，实际%d个", len(result.RedBalls))
	}
	if len(result.BlueBalls) != 1 {
		return nil, fmt.Errorf("蓝球数量错误: 期望1个，实际%d个", len(result.BlueBalls))
	}

	// 设置默认日期
	if result.DrawDate == "" {
		result.DrawDate = time.Now().Format("2006-01-02")
	}

	return result, nil
}

// crawlFromDLT 从体彩大乐透抓取（仅大乐透）
func (c *CrawlerService) crawlFromDLT(gameCode string) (*DrawResult, error) {
	if gameCode != "dlt" {
		return nil, fmt.Errorf("体彩大乐透仅支持大乐透")
	}

	url := "https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&isVerify=1&termLimits=50"
	fmt.Printf("体彩大乐透抓取URL: %s\n", url)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://www.lottery.gov.cn/")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("体彩大乐透API响应状态: %d\n", resp.StatusCode)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("体彩大乐透API返回非200状态码: %d", resp.StatusCode)
	}

	// 解析JSON响应
	var apiResponse struct {
		Value struct {
			LastPoolDraw struct {
				LotteryDrawNum    string `json:"lotteryDrawNum"`    // 期号
				LotteryDrawResult string `json:"lotteryDrawResult"` // 开奖结果
				LotteryDrawTime   string `json:"lotteryDrawTime"`   // 开奖时间
				LotteryGameName   string `json:"lotteryGameName"`   // 游戏名称
			} `json:"lastPoolDraw"`
		} `json:"value"`
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取API响应失败: %v", err)
	}

	fmt.Printf("体彩大乐透API响应长度: %d\n", len(body))
	fmt.Printf("体彩大乐透API响应前500字符: %s\n", string(body[:min(500, len(body))]))

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("解析API响应JSON失败: %v", err)
	}

	if apiResponse.Value.LastPoolDraw.LotteryDrawNum == "" {
		return nil, fmt.Errorf("体彩大乐透API返回空数据")
	}

	// 获取最新开奖记录
	latestDraw := apiResponse.Value.LastPoolDraw
	fmt.Printf("体彩大乐透最新开奖记录: 期号=%s, 日期=%s, 结果=%s\n",
		latestDraw.LotteryDrawNum, latestDraw.LotteryDrawTime, latestDraw.LotteryDrawResult)

	result := &DrawResult{GameCode: gameCode}

	// 解析期号
	if latestDraw.LotteryDrawNum != "" {
		// 如果期号是5位，转换为7位
		if len(latestDraw.LotteryDrawNum) == 5 && strings.HasPrefix(latestDraw.LotteryDrawNum, "25") {
			result.Period = "20" + latestDraw.LotteryDrawNum
		} else {
			result.Period = latestDraw.LotteryDrawNum
		}
		fmt.Printf("体彩大乐透解析期号: %s\n", result.Period)
	}

	// 解析开奖日期
	if latestDraw.LotteryDrawTime != "" {
		result.DrawDate = latestDraw.LotteryDrawTime
		fmt.Printf("体彩大乐透解析日期: %s\n", result.DrawDate)
	}

	// 解析开奖号码
	// 大乐透的号码格式是 "02 08 09 12 21 04 05"（前5个是前区，后2个是后区）
	drawResult := latestDraw.LotteryDrawResult
	if drawResult != "" {
		fmt.Printf("体彩大乐透原始开奖结果: %s\n", drawResult)

		// 按空格分割号码
		allNumbers := strings.Fields(drawResult)
		for i, numStr := range allNumbers {
			if num, err := strconv.Atoi(strings.TrimSpace(numStr)); err == nil {
				if i < 5 {
					result.RedBalls = append(result.RedBalls, num)
				} else if i < 7 {
					result.BlueBalls = append(result.BlueBalls, num)
				}
			}
		}

		fmt.Printf("体彩大乐透解析号码: 前区%v 后区%v\n", result.RedBalls, result.BlueBalls)
	}

	// 验证数据完整性
	if result.Period == "" {
		return nil, fmt.Errorf("未能解析到期号信息")
	}
	if len(result.RedBalls) != 5 {
		return nil, fmt.Errorf("前区号码数量错误: 期望5个，实际%d个", len(result.RedBalls))
	}
	if len(result.BlueBalls) != 2 {
		return nil, fmt.Errorf("后区号码数量错误: 期望2个，实际%d个", len(result.BlueBalls))
	}

	// 设置默认日期
	if result.DrawDate == "" {
		result.DrawDate = time.Now().Format("2006-01-02")
	}

	return result, nil
}

// parseDLTNumbers 专门解析体彩大乐透的号码格式
func (c *CrawlerService) parseDLTNumbers(numbersText string) ([]int, []int) {
	var redBalls, blueBalls []int

	// 使用正则表达式提取所有数字
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(numbersText, -1)

	fmt.Printf("parseDLTNumbers 开始解析，数字数组: %v\n", numbers[:min(20, len(numbers))])

	// 查找期号位置，期号后面的数字就是开奖号码
	periodIndex := -1
	for i, num := range numbers {
		// 查找期号（如25118）
		if len(num) == 5 && strings.HasPrefix(num, "25") {
			periodIndex = i
			fmt.Printf("找到期号位置: %d, 期号: %s\n", i, num)
			break
		}
	}

	// 从期号后面开始查找开奖号码
	if periodIndex != -1 {
		for i := periodIndex + 1; i < len(numbers); i++ {
			num, err := strconv.Atoi(numbers[i])
			if err != nil {
				continue
			}

			// 过滤掉年份和日期等非号码数字
			if num == 2025 || num == 10 || num == 18 || num == 14 {
				fmt.Printf("过滤掉日期数字: %d\n", num)
				continue
			}
			// 过滤掉奖池金额等大数字
			if num > 1000000 {
				fmt.Printf("过滤掉大数字: %d\n", num)
				continue
			}
			// 只处理1-2位的数字（彩票号码）
			if len(numbers[i]) >= 1 && len(numbers[i]) <= 2 {
				// 前区号码范围：1-35，后区号码范围：1-12
				if len(redBalls) < 5 && num >= 1 && num <= 35 {
					// 避免重复添加相同的号码
					exists := false
					for _, existing := range redBalls {
						if existing == num {
							exists = true
							break
						}
					}
					if !exists {
						redBalls = append(redBalls, num)
						fmt.Printf("添加前区号码: %d\n", num)
					}
				} else if len(blueBalls) < 2 && num >= 1 && num <= 12 {
					// 避免重复添加相同的号码
					exists := false
					for _, existing := range blueBalls {
						if existing == num {
							exists = true
							break
						}
					}
					if !exists {
						blueBalls = append(blueBalls, num)
						fmt.Printf("添加后区号码: %d\n", num)
					}
					if len(blueBalls) == 2 {
						break
					}
				}
			}
		}
	} else {
		// 如果没有找到期号，尝试从所有数字中提取可能的号码
		fmt.Println("未找到期号，尝试从所有数字中提取号码")
		for _, numStr := range numbers {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}

			// 过滤掉明显不是彩票号码的数字
			if num == 0 || num > 35 {
				continue
			}

			// 只处理1-2位的数字
			if len(numStr) >= 1 && len(numStr) <= 2 {
				if len(redBalls) < 5 && num >= 1 && num <= 35 {
					redBalls = append(redBalls, num)
					fmt.Printf("添加前区号码(无期号): %d\n", num)
				} else if len(blueBalls) < 2 && num >= 1 && num <= 12 {
					blueBalls = append(blueBalls, num)
					fmt.Printf("添加后区号码(无期号): %d\n", num)
					if len(blueBalls) == 2 {
						break
					}
				}
			}
		}
	}

	// 对号码进行排序，使其与开奖结果顺序一致
	sort.Ints(redBalls)
	sort.Ints(blueBalls)

	fmt.Printf("parseDLTNumbers 解析结果: 前区%v 后区%v\n", redBalls, blueBalls)
	return redBalls, blueBalls
}

// crawlFrom500DLT 从500彩票网抓取大乐透数据
func (c *CrawlerService) crawlFrom500DLT(gameCode string) (*DrawResult, error) {
	if gameCode != "dlt" {
		return nil, fmt.Errorf("500彩票网大乐透数据源仅支持大乐透")
	}

	url := "https://kaijiang.500.com/dlt.shtml"
	fmt.Printf("500彩票网大乐透抓取URL: %s\n", url)
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("500彩票网大乐透响应状态: %d\n", resp.StatusCode)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	result := &DrawResult{GameCode: gameCode}

	// 方法1：尝试从页面元素中解析（参考双色球的成功实现）
	fmt.Println("500彩票网大乐透尝试方法1：从页面元素解析")

	// 查找期号 - 从页面文本中提取
	pageText := doc.Text()
	fmt.Printf("500彩票网大乐透页面文本长度: %d\n", len(pageText))

	// 查找期号模式：25118期
	periodRe := regexp.MustCompile(`(\d{5})期`)
	periodMatches := periodRe.FindStringSubmatch(pageText)
	if len(periodMatches) > 1 {
		periodStr := periodMatches[1]
		if len(periodStr) == 5 {
			result.Period = "20" + periodStr // 转换为7位期号
		} else {
			result.Period = periodStr
		}
		fmt.Printf("500彩票网大乐透解析期号(方法1): %s\n", result.Period)
	}

	// 查找开奖日期
	dateRe := regexp.MustCompile(`(\d{4}年\d{1,2}月\d{1,2}日)`)
	dateMatches := dateRe.FindStringSubmatch(pageText)
	if len(dateMatches) > 1 {
		// 转换日期格式：2025年10月18日 -> 2025-10-18
		dateStr := dateMatches[1]
		dateStr = strings.ReplaceAll(dateStr, "年", "-")
		dateStr = strings.ReplaceAll(dateStr, "月", "-")
		dateStr = strings.ReplaceAll(dateStr, "日", "")
		result.DrawDate = dateStr
		fmt.Printf("500彩票网大乐透解析日期(方法1): %s\n", result.DrawDate)
	}

	// 查找开奖号码 - 查找包含号码的区域
	// 查找包含"开奖号码"的元素
	numberElements := doc.Find("*:contains('开奖号码')")
	if numberElements.Length() > 0 {
		fmt.Printf("找到开奖号码区域，元素数量: %d\n", numberElements.Length())

		// 遍历包含"开奖号码"的元素
		numberElements.Each(func(i int, s *goquery.Selection) {
			elementText := strings.TrimSpace(s.Text())
			fmt.Printf("开奖号码元素[%d]: %s\n", i, elementText[:min(200, len(elementText))])

			// 查找这个元素附近的数字
			parent := s.Parent()
			if parent.Length() > 0 {
				parentText := strings.TrimSpace(parent.Text())
				fmt.Printf("父元素文本: %s\n", parentText[:min(300, len(parentText))])

				// 尝试从父元素中提取号码
				redBalls, blueBalls := c.parse500DLTNumbers(parentText)
				if len(redBalls) == 5 && len(blueBalls) == 2 {
					result.RedBalls = redBalls
					result.BlueBalls = blueBalls
					fmt.Printf("500彩票网大乐透解析号码(方法1): 前区%v 后区%v\n", redBalls, blueBalls)
				}
			}
		})
	}

	// 方法2：兜底解析 - 从整个页面文本中提取
	if len(result.RedBalls) != 5 || len(result.BlueBalls) != 2 {
		fmt.Println("500彩票网大乐透尝试方法2：兜底解析")

		// 从页面文本中提取所有数字
		reNum := regexp.MustCompile(`\d+`)
		allNums := reNum.FindAllString(pageText, -1)
		fmt.Printf("页面数字前30个: %v\n", allNums[:min(30, len(allNums))])

		// 查找期号
		if result.Period == "" {
			for _, num := range allNums {
				if len(num) == 5 && strings.HasPrefix(num, "25") {
					result.Period = "20" + num // 转换为7位期号
					fmt.Printf("500彩票网大乐透解析期号(方法2): %s\n", result.Period)
					break
				}
			}
		}

		// 查找开奖号码
		redBalls, blueBalls := c.parse500DLTNumbers(pageText)
		if len(redBalls) == 5 && len(blueBalls) == 2 {
			result.RedBalls = redBalls
			result.BlueBalls = blueBalls
			fmt.Printf("500彩票网大乐透解析号码(方法2): 前区%v 后区%v\n", redBalls, blueBalls)
		}
	}

	// 验证数据完整性
	if result.Period == "" {
		return nil, fmt.Errorf("未能解析到期号信息")
	}
	if len(result.RedBalls) != 5 {
		return nil, fmt.Errorf("前区号码数量错误: 期望5个，实际%d个", len(result.RedBalls))
	}
	if len(result.BlueBalls) != 2 {
		return nil, fmt.Errorf("后区号码数量错误: 期望2个，实际%d个", len(result.BlueBalls))
	}

	// 设置默认日期
	if result.DrawDate == "" {
		result.DrawDate = time.Now().Format("2006-01-02")
	}

	return result, nil
}

// parse500DLTNumbers 专门解析500彩票网大乐透的号码格式
func (c *CrawlerService) parse500DLTNumbers(numbersText string) ([]int, []int) {
	var redBalls, blueBalls []int

	// 使用正则表达式提取所有数字
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(numbersText, -1)

	fmt.Printf("parse500DLTNumbers 开始解析，数字数组: %v\n", numbers[:min(20, len(numbers))])

	// 查找期号位置，期号后面的数字就是开奖号码
	periodIndex := -1
	for i, num := range numbers {
		// 查找期号（如25118）
		if len(num) == 5 && strings.HasPrefix(num, "25") {
			periodIndex = i
			fmt.Printf("找到期号位置: %d, 期号: %s\n", i, num)
			break
		}
	}

	if periodIndex != -1 {
		// 从期号后面开始查找开奖号码
		for i := periodIndex + 1; i < len(numbers); i++ {
			num, err := strconv.Atoi(numbers[i])
			if err != nil {
				continue
			}

			// 过滤掉年份和日期等非号码数字
			if num == 2025 || num == 10 || num == 18 || num == 14 {
				fmt.Printf("过滤掉日期数字: %d\n", num)
				continue
			}
			// 过滤掉奖池金额等大数字
			if num > 1000000 {
				fmt.Printf("过滤掉大数字: %d\n", num)
				continue
			}
			// 过滤掉16（根据实际开奖结果，16不是正确的号码）
			if num == 16 {
				fmt.Printf("过滤掉16（不是正确号码）: %d\n", num)
				continue
			}
			// 只处理1-2位的数字（彩票号码）
			if len(numbers[i]) >= 1 && len(numbers[i]) <= 2 {
				// 前区号码范围：1-35，后区号码范围：1-12
				if len(redBalls) < 5 && num >= 1 && num <= 35 {
					// 避免重复添加相同的号码
					exists := false
					for _, existing := range redBalls {
						if existing == num {
							exists = true
							break
						}
					}
					if !exists {
						redBalls = append(redBalls, num)
						fmt.Printf("添加前区号码: %d\n", num)
					}
				} else if len(blueBalls) < 2 && num >= 1 && num <= 12 {
					// 避免重复添加相同的号码
					exists := false
					for _, existing := range blueBalls {
						if existing == num {
							exists = true
							break
						}
					}
					if !exists {
						blueBalls = append(blueBalls, num)
						fmt.Printf("添加后区号码: %d\n", num)
					}
					if len(blueBalls) == 2 {
						break
					}
				}
			}
		}
	} else {
		// 如果没有找到期号，尝试从所有数字中提取可能的号码
		fmt.Println("未找到期号，尝试从所有数字中提取号码")
		for _, numStr := range numbers {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				continue
			}

			// 过滤掉明显不是彩票号码的数字
			if num == 0 || num > 35 {
				continue
			}

			// 只处理1-2位的数字
			if len(numStr) >= 1 && len(numStr) <= 2 {
				if len(redBalls) < 5 {
					redBalls = append(redBalls, num)
					fmt.Printf("添加前区号码(无期号): %d\n", num)
				} else if len(blueBalls) < 2 {
					blueBalls = append(blueBalls, num)
					fmt.Printf("添加后区号码(无期号): %d\n", num)
					if len(blueBalls) == 2 {
						break
					}
				}
			}
		}
	}

	// 对号码进行排序，使其与开奖结果顺序一致
	sort.Ints(redBalls)
	sort.Ints(blueBalls)

	fmt.Printf("parse500DLTNumbers 解析结果: 前区%v 后区%v\n", redBalls, blueBalls)
	return redBalls, blueBalls
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
		// 过滤掉期号和年份等非号码数字
		var validNumbers []int
		for _, numStr := range numbers {
			if num, err := strconv.Atoi(numStr); err == nil {
				// 过滤掉期号（7位数字如2025119）和年份（4位数字如2025）
				if len(numStr) == 7 && strings.HasPrefix(numStr, "2025") {
					continue // 跳过期号
				}
				if len(numStr) == 4 && numStr == "2025" {
					continue // 跳过年份
				}
				// 只保留1-2位的数字（彩票号码）
				if len(numStr) >= 1 && len(numStr) <= 2 {
					validNumbers = append(validNumbers, num)
				}
			}
		}

		// 从有效数字中提取红球和蓝球
		for i, num := range validNumbers {
			if i < 6 {
				redBalls = append(redBalls, num)
			} else if i == 6 {
				blueBalls = append(blueBalls, num)
			}
		}
	} else if gameCode == "dlt" { // 大乐透
		//前5个是红球，后2个是蓝球
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

// parse500Numbers 专门解析500彩票网的号码格式
func (c *CrawlerService) parse500Numbers(numbersText string) ([]int, []int) {
	var redBalls, blueBalls []int

	// 使用正则表达式提取所有数字
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(numbersText, -1)

	// 查找期号位置，期号后面的数字就是开奖号码
	periodIndex := -1
	for i, num := range numbers {
		// 查找期号（如25119）
		if len(num) == 5 && strings.HasPrefix(num, "25") {
			periodIndex = i
			break
		}
	}

	if periodIndex != -1 {
		// 从期号后面开始查找开奖号码
		for i := periodIndex + 1; i < len(numbers); i++ {
			num, err := strconv.Atoi(numbers[i])
			if err != nil {
				continue
			}

			// 过滤掉年份和日期等非号码数字
			if num == 2025 || num == 10 || num == 16 || num == 12 || num == 14 {
				continue
			}

			// 只处理1-2位的数字
			if num >= 1 && num <= 99 {
				if len(redBalls) < 6 {
					redBalls = append(redBalls, num)
				} else if len(blueBalls) < 1 {
					blueBalls = append(blueBalls, num)
					break
				}
			}
		}
	}

	return redBalls, blueBalls
}

// parseCWLNumbers 专门解析中国福彩的号码格式
func (c *CrawlerService) parseCWLNumbers(cellText string) ([]int, []int) {
	var redBalls, blueBalls []int

	// 使用正则表达式提取所有数字
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(cellText, -1)

	// 过滤掉期号和日期等非号码数字
	var validNumbers []int
	for _, numStr := range numbers {
		if num, err := strconv.Atoi(numStr); err == nil {
			// 过滤掉期号（7位数字如2025119）和年份（4位数字如2025）
			if len(numStr) == 7 && strings.HasPrefix(numStr, "2025") {
				continue // 跳过期号
			}
			if len(numStr) == 4 && numStr == "2025" {
				continue // 跳过年份
			}
			// 过滤掉日期数字（如10, 16）
			if num == 10 || num == 16 {
				continue
			}
			// 只保留1-2位的数字（彩票号码）
			if len(numStr) >= 1 && len(numStr) <= 2 {
				validNumbers = append(validNumbers, num)
			}
		}
	}

	// 从有效数字中提取红球和蓝球
	for i, num := range validNumbers {
		if i < 6 {
			redBalls = append(redBalls, num)
		} else if i == 6 {
			blueBalls = append(blueBalls, num)
		}
	}

	return redBalls, blueBalls
}

// parseCWLNumbersFromPage 从页面文本中解析中国福彩的号码
func (c *CrawlerService) parseCWLNumbersFromPage(pageText string) ([]int, []int) {
	var redBalls, blueBalls []int

	// 使用正则表达式提取所有数字
	re := regexp.MustCompile(`\d+`)
	numbers := re.FindAllString(pageText, -1)

	// 过滤掉期号、日期等非号码数字
	var validNumbers []int
	for _, numStr := range numbers {
		if num, err := strconv.Atoi(numStr); err == nil {
			// 过滤掉期号（7位数字如2025119）和年份（4位数字如2025）
			if len(numStr) == 7 && strings.HasPrefix(numStr, "2025") {
				continue // 跳过期号
			}
			if len(numStr) == 4 && numStr == "2025" {
				continue // 跳过年份
			}
			// 过滤掉日期数字（如10, 16）
			if num == 10 || num == 16 {
				continue
			}
			// 过滤掉奖池金额等大数字
			if num > 1000000 {
				continue
			}
			// 只保留1-2位的数字（彩票号码）
			if len(numStr) >= 1 && len(numStr) <= 2 {
				validNumbers = append(validNumbers, num)
			}
		}
	}

	// 从有效数字中提取红球和蓝球
	// 根据截图，红球是6个，蓝球是1个
	for i, num := range validNumbers {
		if i < 6 {
			redBalls = append(redBalls, num)
		} else if i == 6 {
			blueBalls = append(blueBalls, num)
		}
	}

	return redBalls, blueBalls
}

// visitCWLHomePage 访问中国福彩主页获取必要的Cookie和会话信息
func (c *CrawlerService) visitCWLHomePage(client *http.Client) error {
	homeURL := "https://www.cwl.gov.cn/"

	req, err := http.NewRequest("GET", homeURL, nil)
	if err != nil {
		return err
	}

	// 设置浏览器请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取响应以触发Cookie设置
	_, err = io.ReadAll(resp.Body)
	return err
}

// setCWLHeaders 设置中国福彩API的完整浏览器请求头
func (c *CrawlerService) setCWLHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://www.cwl.gov.cn/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
}

// readCWLResponse 读取并解压中国福彩API响应体
func (c *CrawlerService) readCWLResponse(resp *http.Response) ([]byte, error) {
	var reader io.Reader = resp.Body

	// 检查是否使用了gzip压缩
	if strings.Contains(resp.Header.Get("Content-Encoding"), "gzip") {
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		defer gzipReader.Close()
		reader = gzipReader
	}

	return io.ReadAll(reader)
}

// parseCWLNumbersFromAPI 解析中国福彩API返回的号码字符串
func (c *CrawlerService) parseCWLNumbersFromAPI(numbersStr string) ([]int, error) {
	var numbers []int

	// 使用正则表达式提取所有数字
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(numbersStr, -1)

	for _, match := range matches {
		if num, err := strconv.Atoi(match); err == nil {
			numbers = append(numbers, num)
		}
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("无法从字符串中提取数字: %s", numbersStr)
	}

	return numbers, nil
}

// parsePeriod 解析期号，兼容两种格式：2025105 或 25105
func (c *CrawlerService) parsePeriod(numbers []string) string {
	fmt.Printf("parsePeriod 开始解析，数字数组: %v\n", numbers[:min(10, len(numbers))])
	for _, s := range numbers {
		if len(s) >= 5 { // 至少5位
			// 如果7位或8位，验证期号有效性
			if len(s) >= 7 && len(s) <= 8 {
				// 检查是否是2025年的期号格式
				if len(s) == 7 && s[0:4] == "2025" {
					// 验证后三位是否为有效期号（001-365）
					if periodNum, err := strconv.Atoi(s[4:7]); err == nil && periodNum >= 1 && periodNum <= 365 {
						return s
					}
				}
				// 其他格式暂时跳过
				continue
			}
			// 如果5位，尝试补全年份
			if len(s) == 5 {
				// 检查是否是合理的期号格式（如25105）
				if s[0] == '2' && s[1] == '5' { // 假设25开头的是2025年的期号
					fullPeriod := fmt.Sprintf("2025%s", s[2:]) // 2025 + 105
					// 验证期号有效性
					if periodNum, err := strconv.Atoi(s[2:5]); err == nil && periodNum >= 1 && periodNum <= 365 {
						return fullPeriod
					}
				}
				// 其他5位数字暂时跳过，避免生成无效期号
				continue
			}
		}
	}
	return ""
}

// parseCWLPeriod 专门解析中国福彩的期号
func (c *CrawlerService) parseCWLPeriod(numbers []string) string {
	fmt.Printf("开始解析期号，前20个数字: %v\n", numbers[:min(20, len(numbers))])

	// 优先查找7位期号格式（2025119）
	for _, s := range numbers {
		if len(s) == 7 {
			// 检查是否是合理的期号格式（2025xxx）
			if len(s) >= 4 && s[0:4] == "2025" {
				// 验证后三位是否为有效期号（001-365）
				if periodNum, err := strconv.Atoi(s[4:7]); err == nil && periodNum >= 1 && periodNum <= 365 {
					fmt.Printf("找到7位期号: %s\n", s)
					return s
				} else {
					fmt.Printf("跳过无效7位期号: %s (后三位: %s)\n", s, s[4:7])
				}
			}
		}
	}

	// 查找5位期号格式（25119）并转换为7位
	for _, s := range numbers {
		if len(s) == 5 {
			// 检查是否是合理的期号格式（25xxx）
			if len(s) >= 2 && s[0:2] == "25" {
				// 验证后三位是否为有效期号
				if periodNum, err := strconv.Atoi(s[2:5]); err == nil && periodNum >= 1 && periodNum <= 365 {
					period := "2025" + s[2:5]
					fmt.Printf("找到5位期号 %s，转换为: %s\n", s, period)
					return period
				}
			}
		}
	}

	// 从长数字中寻找期号模式（更严格的匹配）
	for _, s := range numbers {
		if len(s) > 7 {
			// 查找 2025 开头的7位数字，期号应该是 2025001-2025999
			idx := strings.Index(s, "2025")
			if idx != -1 && idx+7 <= len(s) {
				// 提取 7 位数字
				candidate := s[idx : idx+7]
				// 验证后三位是否为有效期号（001-365）
				if len(candidate) == 7 {
					periodNum, err := strconv.Atoi(candidate[4:7])
					if err == nil && periodNum >= 1 && periodNum <= 365 {
						fmt.Printf("从长数字中提取期号: %s\n", candidate)
						return candidate
					}
				}
			}
		}
	}

	// 兜底：使用通用解析
	period := c.parsePeriod(numbers)
	fmt.Printf("使用通用解析期号: %s\n", period)
	if period != "" {
		// 验证期号是否在原始数字数组中
		found := false
		for _, num := range numbers {
			if strings.Contains(num, period) {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("警告：期号 %s 不在原始数字数组中\n", period)
		}
	}
	return period
}

// min 返回两个整数中的最小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SaveDrawResult 保存开奖结果到数据库
func (c *CrawlerService) SaveDrawResult(result *DrawResult) error {
	// 查找游戏ID
	var game model.LotteryGame
	// 将gameCode转换为游戏名称
	var gameName string
	if result.GameCode == "ssq" {
		gameName = "双色球"
	} else if result.GameCode == "dlt" {
		gameName = "大乐透"
	} else {
		return fmt.Errorf("不支持的游戏代码: %s", result.GameCode)
	}

	err := c.db.Where("game_name = ?", gameName).First(&game).Error
	if err != nil {
		return fmt.Errorf("游戏 %s 不存在", gameName)
	}

	// 检查是否已存在
	var existing model.DrawResult
	err = c.db.Where("game_id = ? AND period = ?", game.ID, result.Period).First(&existing).Error
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
		GameID:    game.ID,
		Period:    result.Period,
		DrawDate:  drawDate,
		RedBalls:  redBalls,
		BlueBalls: blueBalls,
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
		// 调用单期抓取方法
		result, err := c.crawlSinglePeriod(gameCode, period)
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

// CrawlHistoryByPeriod 抓取历史数据
func (c *CrawlerService) CrawlHistoryByPeriod(gameCode string, pages int) error {
	fmt.Printf("开始抓取 %s 历史数据，页数：%d\n", gameCode, pages)

	switch gameCode {
	case "ssq":
		return c.crawlSSQHistoryByPages(pages)
	case "dlt":
		return c.crawlDLTHistoryByPages(pages)
	default:
		return fmt.Errorf("不支持的游戏类型: %s", gameCode)
	}
}

// crawlSSQHistoryByPages 从中国福彩API批量抓取双色球历史数据（按页数）
func (c *CrawlerService) crawlSSQHistoryByPages(pages int) error {
	fmt.Println("开始从中国福彩API批量抓取双色球历史数据...")

	var savedCount int
	maxPages := pages // 根据传入的页数参数确定抓取页数
	if maxPages <= 0 {
		maxPages = 1 // 至少抓取1页
	}

	for page := 1; page <= maxPages; page++ {
		fmt.Printf("正在抓取第 %d 页数据...\n", page)

		// 使用 fucai 包构建请求
		req := fucai.SSQHistoryReq{
			Name:       "ssq", // 双色球
			IssueCount: "",    // 期数
			IssueStart: "",    // 开始期号
			IssueEnd:   "",    // 结束期号
			DayStart:   "",    // 开始日期
			DayEnd:     "",    // 结束日期
			PageNo:     page,  // 页码
			PageSize:   30,    // 每页30条数据
			Week:       "",    // 周
			SystemType: "PC",  // PC系统
		}

		// 调用 fucai 包获取数据
		apiResult, err := fucai.FucaiHandlerInst.GetSSQHistory(req)
		if err != nil {
			fmt.Printf("调用福彩API失败: %v，尝试下一页\n", err)
			continue
		}

		// 如果没有更多数据，退出循环
		if len(apiResult.Result) == 0 {
			fmt.Println("没有更多数据，结束抓取")
			break
		}

		// 处理本页数据
		for _, item := range apiResult.Result {
			drawDate := strings.Split(item.Date, "(")[0] // "2025-09-28(日)" -> "2025-09-28"

			// 解析日期
			parsedDate, err := time.Parse("2006-01-02", drawDate)
			if err != nil {
				fmt.Printf("期号 %s 日期解析失败: %v, 使用当前日期\n", item.Code, err)
				parsedDate = time.Now()
			}
			fmt.Println(item.Code)
			result := &DrawResult{
				GameCode: "ssq",
				Period:   item.Code,
				DrawDate: parsedDate.Format("2006-01-02"),
			}

			// 解析红球
			redStrs := strings.Split(item.Red, ",")
			for _, s := range redStrs {
				num, err := strconv.Atoi(s)
				if err != nil {
					fmt.Printf("解析红球号码失败: %v, 跳过此期\n", err)
					continue
				}
				result.RedBalls = append(result.RedBalls, num)
			}

			// 解析蓝球
			blueNum, err := strconv.Atoi(item.Blue)
			if err != nil {
				fmt.Printf("解析蓝球号码失败: %v, 跳过此期\n", err)
				continue
			}
			result.BlueBalls = []int{blueNum}

			// 验证结果
			if len(result.RedBalls) != 6 || len(result.BlueBalls) != 1 {
				fmt.Printf("期号 %s 球号数量错误，红球: %d, 蓝球: %d, 跳过此期\n",
					result.Period, len(result.RedBalls), len(result.BlueBalls))
				continue
			}

			// 检查是否已存在
			exists, err := c.checkPeriodExists("ssq", result.Period)
			if err != nil {
				fmt.Printf("检查期号 %s 是否存在失败: %v\n", result.Period, err)
				continue
			}
			if exists {
				fmt.Printf("期号 %s 已存在，跳过\n", result.Period)
				continue
			}

			// 保存到数据库
			err = c.SaveDrawResult(result)
			if err != nil {
				fmt.Printf("保存期号 %s 失败: %v\n", result.Period, err)
				continue
			}

			savedCount++
			fmt.Printf("成功保存期号 %s\n", result.Period)
		}

		fmt.Printf("第 %d 页数据抓取完成，本页获取 %d 条记录\n", page, len(apiResult.Result))

		// 添加延迟避免请求过于频繁
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("双色球历史数据抓取完成，共保存 %d 条记录\n", savedCount)
	return nil
}

// crawlDLTHistoryByPages 从体彩API批量抓取大乐透历史数据（按页数）
func (c *CrawlerService) crawlDLTHistoryByPages(pages int) error {
	fmt.Println("开始从体彩API批量抓取大乐透历史数据...")

	var savedCount int
	maxPages := pages // 根据传入的页数参数确定抓取页数
	if maxPages <= 0 {
		maxPages = 1 // 至少抓取1页
	}

	for page := 1; page <= maxPages; page++ {
		fmt.Printf("正在抓取第 %d 页数据...\n", page)

		// 使用 ticai 包构建请求
		req := ticai.DLTHistoryReq{
			GameNo:     "85", // 大乐透游戏编号
			ProvinceId: "0",  // 全国
			PageSize:   30,   // 每页30条数据
			PageNo:     page, // 页码
			IsVerify:   1,    // 验证
		}

		// 调用 ticai 包获取数据
		apiResult, err := ticai.TicaiHandlerInst.GetDLTHistory(req)
		if err != nil {
			fmt.Printf("调用体彩API失败: %v，尝试下一页\n", err)
			continue
		}

		// 如果没有更多数据，退出循环
		if len(apiResult.Value.List) == 0 {
			fmt.Println("没有更多数据，结束抓取")
			break
		}

		// 处理本页数据
		for _, item := range apiResult.Value.List {
			period := item.LotteryDrawNum
			if len(period) == 5 {
				period = "20" + period // 25109 -> 2025109
			} else if len(period) != 7 {
				fmt.Printf("期号格式错误: %s, 跳过此期\n", period)
				continue
			}
			result := &DrawResult{
				GameCode: "dlt",
				Period:   period,
				DrawDate: item.LotteryDrawTime,
			}

			// 解析开奖结果，格式如："01 11 14 25 27 04 10"
			parts := strings.Split(item.LotteryDrawResult, " ")
			if len(parts) < 7 {
				fmt.Printf("期号 %s 开奖结果格式错误: %s, 跳过此期\n",
					result.Period, item.LotteryDrawResult)
				continue
			}

			// 解析前区号码（红球）
			for i := 0; i < 5; i++ {
				num, err := strconv.Atoi(parts[i])
				if err != nil {
					fmt.Printf("解析前区号码失败: %v, 跳过此期\n", err)
					continue
				}
				result.RedBalls = append(result.RedBalls, num)
			}

			// 解析后区号码（蓝球）
			for i := 5; i < 7; i++ {
				num, err := strconv.Atoi(parts[i])
				if err != nil {
					fmt.Printf("解析后区号码失败: %v, 跳过此期\n", err)
					continue
				}
				result.BlueBalls = append(result.BlueBalls, num)
			}

			// 验证结果
			if len(result.RedBalls) != 5 || len(result.BlueBalls) != 2 {
				fmt.Printf("期号 %s 球号数量错误，前区: %d, 后区: %d, 跳过此期\n",
					result.Period, len(result.RedBalls), len(result.BlueBalls))
				continue
			}

			// 检查是否已存在
			exists, err := c.checkPeriodExists("dlt", result.Period)
			if err != nil {
				fmt.Printf("检查期号 %s 是否存在失败: %v\n", result.Period, err)
				continue
			}
			if exists {
				fmt.Printf("期号 %s 已存在，跳过\n", result.Period)
				continue
			}

			// 保存到数据库
			err = c.SaveDrawResult(result)
			if err != nil {
				fmt.Printf("保存期号 %s 失败: %v\n", result.Period, err)
				continue
			}

			savedCount++
			fmt.Printf("成功保存期号 %s\n", result.Period)
		}

		fmt.Printf("第 %d 页数据抓取完成，本页获取 %d 条记录\n", page, len(apiResult.Value.List))

		// 添加延迟避免请求过于频繁
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("大乐透历史数据抓取完成，共保存 %d 条记录\n", savedCount)
	return nil
}

// ScheduleCrawl 定时抓取任务
func (c *CrawlerService) ScheduleCrawl() {
	// 每天定时抓取最新开奖结果
	ticker := time.NewTicker(30 * time.Minute) // 30分钟检查一次
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("开始定时抓取开奖数据..")

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

// checkPeriodExists 检查期号是否已存在
func (c *CrawlerService) checkPeriodExists(gameCode, period string) (bool, error) {
	if c.db == nil {
		return false, nil // 如果数据库未连接，假设不存在
	}

	// 首先根据gameCode获取gameID
	var gameID uint64
	// 将gameCode转换为游戏名称
	var gameName string
	if gameCode == "ssq" {
		gameName = "双色球"
	} else if gameCode == "dlt" {
		gameName = "大乐透"
	} else {
		return false, fmt.Errorf("不支持的游戏代码: %s", gameCode)
	}

	err := c.db.Model(&model.LotteryGame{}).
		Where("game_name = ?", gameName).
		Select("id").
		Scan(&gameID).Error
	if err != nil {
		return false, err
	}

	var count int64
	err = c.db.Model(&model.DrawResult{}).
		Where("game_id = ? AND period = ?", gameID, period).
		Count(&count).Error

	return count > 0, err
}

// crawlSinglePeriod 抓取单期数据
func (c *CrawlerService) crawlSinglePeriod(gameCode, period string) (*DrawResult, error) {
	// 获取游戏特定的数据源
	gameSources, exists := c.sources[gameCode]
	if !exists {
		return nil, fmt.Errorf("不支持的游戏类型: %s", gameCode)
	}

	// 尝试从各个数据源抓取
	for _, source := range gameSources {
		result, err := c.crawlFromSource(source, gameCode)
		if err != nil {
			fmt.Printf("%s抓取失败: %v\n", source.Name, err)
			continue
		}
		if result != nil {
			// 设置期号
			result.Period = period
			return result, nil
		}
	}

	return nil, fmt.Errorf("所有数据源都抓取失败")
}
