package fucai

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

// GetSSQHistory 获取双色球历史数据
func (h *FucaiHandler) GetSSQHistory(req SSQHistoryReq) (res SSQHistoryResp, err error) {
	res = SSQHistoryResp{}

	// 创建Cookie Jar来管理会话
	jar, err := cookiejar.New(nil)
	if err != nil {
		return res, fmt.Errorf("创建Cookie Jar失败: %v", err)
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
	}

	// 第一步：先访问主页获取Cookie和会话信息
	err = h.visitHomePage(client)
	if err != nil {
		fmt.Printf("访问主页失败: %v\n", err)
		// 不要因为这个失败就退出，继续尝试
	}

	// 添加随机延迟，模拟真实用户行为
	time.Sleep(2 * time.Second)

	// 第二步：发送API请求
	url := fmt.Sprintf("https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=%s&issueCount=%s&issueStart=%s&issueEnd=%s&dayStart=%s&dayEnd=%s&pageNo=%d&pageSize=%d&week=%s&systemType=%s",
		req.Name, req.IssueCount, req.IssueStart, req.IssueEnd, req.DayStart, req.DayEnd, req.PageNo, req.PageSize, req.Week, req.SystemType)

	fmt.Println("请求URL: ", url)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置完整的浏览器请求头
	h.setHeaders(httpReq)

	// 发送请求
	resp, err := client.Do(httpReq)
	if err != nil {
		return res, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return res, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取并解压响应体
	responseBytes, err := h.readResponse(resp)
	if err != nil {
		return res, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON响应
	err = json.Unmarshal(responseBytes, &res)
	if err != nil {
		return res, fmt.Errorf("解析JSON失败: %v", err)
	}

	// 检查API错误码
	if res.State != 0 {
		return res, fmt.Errorf("API返回错误: %s", res.Message)
	}

	return res, nil
}

// visitHomePage 访问主页获取必要的Cookie和会话信息
func (h *FucaiHandler) visitHomePage(client *http.Client) error {
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

// setHeaders 设置完整的浏览器请求头
func (h *FucaiHandler) setHeaders(req *http.Request) {
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

// readResponse 读取并解压响应体
func (h *FucaiHandler) readResponse(resp *http.Response) ([]byte, error) {
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
