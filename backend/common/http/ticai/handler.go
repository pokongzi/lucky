package ticai

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GetDLTHistory 获取大乐透历史数据
func (h *TicaiHandler) GetDLTHistory(req DLTHistoryReq) (res DLTHistoryResp, err error) {
	res = DLTHistoryResp{}

	// 构建请求URL
	url := fmt.Sprintf("https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=%s&provinceId=%s&pageSize=%d&isVerify=%d&pageNo=%d",
		req.GameNo, req.ProvinceId, req.PageSize, req.IsVerify, req.PageNo)

	// 创建HTTP客户端
	client := &http.Client{Timeout: 10 * time.Second}

	// 创建请求
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return res, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头，模拟真实浏览器
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	httpReq.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	httpReq.Header.Set("Accept-Encoding", "gzip, deflate, br")
	httpReq.Header.Set("Referer", "https://webapi.sporttery.cn/")
	httpReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	httpReq.Header.Set("Connection", "keep-alive")
	httpReq.Header.Set("Sec-Fetch-Dest", "empty")
	httpReq.Header.Set("Sec-Fetch-Mode", "cors")
	httpReq.Header.Set("Sec-Fetch-Site", "same-origin")
	httpReq.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
	httpReq.Header.Set("sec-ch-ua-mobile", "?0")
	httpReq.Header.Set("sec-ch-ua-platform", `"Windows"`)

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
	if res.ErrorCode != "0" {
		return res, fmt.Errorf("API返回错误: %s", res.ErrorMessage)
	}

	// 处理开奖结果，拆分为红球和蓝球
	for i := range res.Value.List {
		// 拆分开奖结果，格式为：05 07 08 15 33 06 10
		// 大乐透规则：前5个数字为红球，后2个数字为蓝球
		parts := strings.Split(res.Value.List[i].LotteryDrawResult, " ")
		if len(parts) >= 7 { // 确保至少有7个数字
			// 前5个数字为红球
			redBalls := strings.Join(parts[:5], " ")
			// 后2个数字为蓝球
			blueBalls := strings.Join(parts[5:7], " ")

			res.Value.List[i].RedBalls = redBalls   // 红球部分
			res.Value.List[i].BlueBalls = blueBalls // 蓝球部分
		}
	}

	return res, nil
}

// readResponse 读取并解压响应体
func (h *TicaiHandler) readResponse(resp *http.Response) ([]byte, error) {
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
