package fucai

import (
	"encoding/json"
	"testing"
)

// TestGetSSQHistory_Success 测试获取双色球历史数据
func TestGetSSQHistory_Success(t *testing.T) {
	// 创建处理器实例
	handler := &FucaiHandler{}

	// 创建请求参数
	req := SSQHistoryReq{
		Name:       "ssq", // 双色球
		IssueCount: "",    // 期数
		IssueStart: "",    // 开始期号
		IssueEnd:   "",    // 结束期号
		DayStart:   "",    // 开始日期
		DayEnd:     "",    // 结束日期
		PageNo:     1,     // 页码
		PageSize:   5,     // 设置较小的页面大小以减少输出
		Week:       "",    // 周
		SystemType: "PC",  // PC系统
	}

	// 调用方法
	resp, err := handler.GetSSQHistory(req)

	// 打印结果
	if err != nil {
		t.Logf("GetSSQHistory returned error: %v", err)
		return
	}

	// 打印响应状态
	t.Logf("API Response State: %d", resp.State)
	t.Logf("API Response Message: %s", resp.Message)
	t.Logf("API Response Total: %d", resp.Total)
	t.Logf("API Response TFoot: %s", resp.TFoot)
	t.Logf("API Response Result Count: %d", len(resp.Result))

	// 打印前几条记录
	for i, item := range resp.Result {
		if i >= 3 { // 只打印前3条
			break
		}
		t.Logf("Record %d - Code: %s, Date: %s, Red: %s, Blue: %s",
			i+1, item.Code, item.Date, item.Red, item.Blue)
	}

	// 打印完整响应的JSON格式（用于调试）
	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	t.Logf("Complete Response JSON:\n%s", string(respJSON))
}
