package http500

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// SSQBlueMissingData 双色球蓝球遗漏数据结构
type SSQBlueMissingData struct {
	Number         int     `json:"number"`         // 号码
	Theoretical    float64 `json:"theoretical"`    // 理论次数
	Count          int     `json:"count"`          // 出现次数
	LastMissing    int     `json:"lastMissing"`    // 上次遗漏
	CurrentMissing int     `json:"currentMissing"` // 本次遗漏
	MaxMissing     int     `json:"maxMissing"`     // 最大遗漏
}

// FetchSSQBlueMissingData 抓取双色球蓝球遗漏数据
// periodCount: 期数范围，支持10/30/50
func FetchSSQBlueMissingData(periodCount int) ([]SSQBlueMissingData, error) {
	// 确保期数参数有效
	if periodCount != 10 && periodCount != 30 && periodCount != 50 {
		periodCount = 30 // 默认30期
	}

	// 构建请求URL
	url := fmt.Sprintf("https://datachart.500.com/ssq/omit/newinc/hmyl_blue.php?select=%d", periodCount)

	// 创建HTTP客户端并设置超时时间
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Referer", "https://datachart.500.com/ssq/")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}

	// 使用goquery解析HTML
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return nil, fmt.Errorf("解析HTML失败: %w", err)
	}

	// 存储结果
	var missingData []SSQBlueMissingData

	// 解析表格数据 - 使用正确的表格选择器
	doc.Find("table.table_list001 tr").Each(func(i int, s *goquery.Selection) {
		// 跳过表头行
		if i == 0 {
			return
		}

		// 提取数据
		var data SSQBlueMissingData
		var cols []string

		s.Find("td").Each(func(j int, td *goquery.Selection) {
			cols = append(cols, strings.TrimSpace(td.Text()))
		})

		// 确保有足够的列
		if len(cols) < 11 { // 根据HTML，这个表格有11列
			return
		}

		// 解析号码 (第1列)
		number, err := strconv.Atoi(cols[0])
		if err != nil {
			return
		}
		data.Number = number

		// 解析理论次数 (第3列) - 直接从接口获取
		theoretical, err := strconv.ParseFloat(cols[2], 64)
		if err != nil {
			theoretical = 0
		}
		data.Theoretical = theoretical

		// 解析出现次数 (第2列)
		count, err := strconv.Atoi(cols[1])
		if err != nil {
			return
		}
		data.Count = count

		// 解析上次遗漏 (第7列)
		lastMissing, err := strconv.Atoi(cols[6])
		if err != nil {
			lastMissing = 0
		}
		data.LastMissing = lastMissing

		// 解析本次遗漏 (第8列)
		currentMissing, err := strconv.Atoi(cols[7])
		if err != nil {
			currentMissing = 0
		}
		data.CurrentMissing = currentMissing

		// 解析最大遗漏 (第6列)
		maxMissing, err := strconv.Atoi(cols[5])
		if err != nil {
			maxMissing = 0
		}
		data.MaxMissing = maxMissing

		// 添加到结果集
		missingData = append(missingData, data)
	})

	// 如果没有数据，返回错误
	if len(missingData) == 0 {
		return nil, fmt.Errorf("未找到遗漏数据")
	}

	return missingData, nil
}

// GetSSQBlueMissingDataJSON 获取双色球蓝球遗漏数据并返回JSON字符串
func GetSSQBlueMissingDataJSON(periodCount int) (string, error) {
	data, err := FetchSSQBlueMissingData(periodCount)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON编码失败: %w", err)
	}

	return string(jsonBytes), nil
}
