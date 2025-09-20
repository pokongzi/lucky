package ticai

import (
	"encoding/json"
	"testing"
)

// TestGetDLTHistory_Success 测试获取大乐透历史数据
func TestGetDLTHistory_Success(t *testing.T) {
	// 创建处理器实例
	handler := &TicaiHandler{}

	// 创建请求参数
	req := DLTHistoryReq{
		GameNo:     "85", // 大乐透游戏编号
		ProvinceId: "0",  // 全国
		PageSize:   5,    // 设置较小的页面大小以减少输出
		PageNo:     1,    // 页码
		IsVerify:   1,    // 验证
	}

	// 调用方法
	resp, err := handler.GetDLTHistory(req)

	// 打印结果
	if err != nil {
		t.Logf("GetDLTHistory returned error: %v", err)
		return
	}

	// 打印响应状态
	t.Logf("API Response ErrorCode: %s", resp.ErrorCode)
	t.Logf("API Response ErrorMessage: %s", resp.ErrorMessage)
	t.Logf("API Response PageNo: %d", resp.Value.PageNo)
	t.Logf("API Response PageSize: %d", resp.Value.PageSize)
	t.Logf("API Response Total: %d", resp.Value.Total)
	t.Logf("API Response List Count: %d", len(resp.Value.List))

	// 打印前几条记录
	for i, item := range resp.Value.List {
		if i >= 3 { // 只打印前3条
			break
		}
		t.Logf("Record %d - DrawNum: %s, DrawTime: %s, DrawResult: %s, RedBalls: %s, BlueBalls: %s",
			i+1, item.LotteryDrawNum, item.LotteryDrawTime, item.LotteryDrawResult, item.RedBalls, item.BlueBalls)
	}

	// 打印完整响应的JSON格式（用于调试）
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Complete Response JSON:\n%s", string(respJSON))
}
