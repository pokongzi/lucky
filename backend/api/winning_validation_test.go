package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSSQWinningRulesComplete 完整测试双色球所有中奖规则
func TestSSQWinningRulesComplete(t *testing.T) {
	tests := []struct {
		name          string
		redMatches    int
		blueMatches   int
		expectedWin   string
		expectedPrize int64
		description   string
	}{
		// 一等奖
		{
			name:          "SSQ 一等奖 (6红+1蓝)",
			redMatches:    6,
			blueMatches:   1,
			expectedWin:   "一等奖",
			expectedPrize: 500000000,
			description:   "选中6个红球+1个蓝球",
		},
		// 二等奖
		{
			name:          "SSQ 二等奖 (6红+0蓝)",
			redMatches:    6,
			blueMatches:   0,
			expectedWin:   "二等奖",
			expectedPrize: 10000000,
			description:   "选中6个红球",
		},
		// 三等奖
		{
			name:          "SSQ 三等奖 (5红+1蓝)",
			redMatches:    5,
			blueMatches:   1,
			expectedWin:   "三等奖",
			expectedPrize: 300000,
			description:   "选中5个红球+1个蓝球",
		},
		// 四等奖
		{
			name:          "SSQ 四等奖 (5红+0蓝)",
			redMatches:    5,
			blueMatches:   0,
			expectedWin:   "四等奖",
			expectedPrize: 20000,
			description:   "选中5个红球",
		},
		{
			name:          "SSQ 四等奖 (4红+1蓝)",
			redMatches:    4,
			blueMatches:   1,
			expectedWin:   "四等奖",
			expectedPrize: 20000,
			description:   "选中4个红球+1个蓝球",
		},
		// 五等奖
		{
			name:          "SSQ 五等奖 (4红+0蓝)",
			redMatches:    4,
			blueMatches:   0,
			expectedWin:   "五等奖",
			expectedPrize: 1000,
			description:   "选中4个红球",
		},
		{
			name:          "SSQ 五等奖 (3红+1蓝)",
			redMatches:    3,
			blueMatches:   1,
			expectedWin:   "五等奖",
			expectedPrize: 1000,
			description:   "选中3个红球+1个蓝球",
		},
		// 六等奖
		{
			name:          "SSQ 六等奖 (2红+1蓝)",
			redMatches:    2,
			blueMatches:   1,
			expectedWin:   "六等奖",
			expectedPrize: 500,
			description:   "选中2个红球+1个蓝球",
		},
		{
			name:          "SSQ 六等奖 (1红+1蓝)",
			redMatches:    1,
			blueMatches:   1,
			expectedWin:   "六等奖",
			expectedPrize: 500,
			description:   "选中1个红球+1个蓝球",
		},
		{
			name:          "SSQ 六等奖 (0红+1蓝)",
			redMatches:    0,
			blueMatches:   1,
			expectedWin:   "六等奖",
			expectedPrize: 500,
			description:   "只选中1个蓝球",
		},
		// 未中奖情况
		{
			name:          "SSQ 未中奖 (3红+0蓝)",
			redMatches:    3,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
			description:   "选中3个红球但没有蓝球",
		},
		{
			name:          "SSQ 未中奖 (2红+0蓝)",
			redMatches:    2,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
			description:   "选中2个红球但没有蓝球",
		},
		{
			name:          "SSQ 未中奖 (1红+0蓝)",
			redMatches:    1,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
			description:   "只选中1个红球",
		},
		{
			name:          "SSQ 未中奖 (0红+0蓝)",
			redMatches:    0,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
			description:   "一个都没选中",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winLevel, prizeAmount := calculateWinLevel("ssq", tt.redMatches, tt.blueMatches)
			assert.Equal(t, tt.expectedWin, winLevel, "中奖等级不匹配: %s", tt.description)
			assert.Equal(t, tt.expectedPrize, prizeAmount, "奖金金额不匹配: %s", tt.description)
			t.Logf("✓ %s - %s", tt.name, tt.description)
		})
	}
}

// TestDLTWinningRulesComplete 完整测试大乐透所有中奖规则
func TestDLTWinningRulesComplete(t *testing.T) {
	tests := []struct {
		name          string
		redMatches    int
		blueMatches   int
		expectedWin   string
		expectedPrize int64
		description   string
	}{
		// 一等奖
		{
			name:          "DLT 一等奖 (5红+2蓝)",
			redMatches:    5,
			blueMatches:   2,
			expectedWin:   "一等奖",
			expectedPrize: 1000000000,
			description:   "选中5个前区+2个后区",
		},
		// 二等奖
		{
			name:          "DLT 二等奖 (5红+1蓝)",
			redMatches:    5,
			blueMatches:   1,
			expectedWin:   "二等奖",
			expectedPrize: 80000000,
			description:   "选中5个前区+1个后区",
		},
		// 三等奖
		{
			name:          "DLT 三等奖 (5红+0蓝)",
			redMatches:    5,
			blueMatches:   0,
			expectedWin:   "三等奖",
			expectedPrize: 1000000,
			description:   "选中5个前区",
		},
		// 四等奖
		{
			name:          "DLT 四等奖 (4红+2蓝)",
			redMatches:    4,
			blueMatches:   2,
			expectedWin:   "四等奖",
			expectedPrize: 300000,
			description:   "选中4个前区+2个后区",
		},
		// 五等奖
		{
			name:          "DLT 五等奖 (4红+1蓝)",
			redMatches:    4,
			blueMatches:   1,
			expectedWin:   "五等奖",
			expectedPrize: 30000,
			description:   "选中4个前区+1个后区",
		},
		{
			name:          "DLT 五等奖 (3红+2蓝)",
			redMatches:    3,
			blueMatches:   2,
			expectedWin:   "五等奖",
			expectedPrize: 30000,
			description:   "选中3个前区+2个后区",
		},
		// 六等奖
		{
			name:          "DLT 六等奖 (4红+0蓝)",
			redMatches:    4,
			blueMatches:   0,
			expectedWin:   "六等奖",
			expectedPrize: 10000,
			description:   "选中4个前区",
		},
		{
			name:          "DLT 六等奖 (3红+1蓝)",
			redMatches:    3,
			blueMatches:   1,
			expectedWin:   "六等奖",
			expectedPrize: 10000,
			description:   "选中3个前区+1个后区",
		},
		{
			name:          "DLT 六等奖 (2红+2蓝)",
			redMatches:    2,
			blueMatches:   2,
			expectedWin:   "六等奖",
			expectedPrize: 10000,
			description:   "选中2个前区+2个后区",
		},
		// 七等奖
		{
			name:          "DLT 七等奖 (3红+0蓝)",
			redMatches:    3,
			blueMatches:   0,
			expectedWin:   "七等奖",
			expectedPrize: 1500,
			description:   "选中3个前区",
		},
		{
			name:          "DLT 七等奖 (2红+1蓝)",
			redMatches:    2,
			blueMatches:   1,
			expectedWin:   "七等奖",
			expectedPrize: 1500,
			description:   "选中2个前区+1个后区",
		},
		{
			name:          "DLT 七等奖 (1红+2蓝)",
			redMatches:    1,
			blueMatches:   2,
			expectedWin:   "七等奖",
			expectedPrize: 1500,
			description:   "选中1个前区+2个后区",
		},
		{
			name:          "DLT 七等奖 (0红+2蓝)",
			redMatches:    0,
			blueMatches:   2,
			expectedWin:   "七等奖",
			expectedPrize: 1500,
			description:   "只选中2个后区",
		},
		// 八等奖
		{
			name:          "DLT 八等奖 (2红+0蓝)",
			redMatches:    2,
			blueMatches:   0,
			expectedWin:   "八等奖",
			expectedPrize: 500,
			description:   "选中2个前区",
		},
		{
			name:          "DLT 八等奖 (1红+1蓝)",
			redMatches:    1,
			blueMatches:   1,
			expectedWin:   "八等奖",
			expectedPrize: 500,
			description:   "选中1个前区+1个后区",
		},
		{
			name:          "DLT 八等奖 (0红+1蓝)",
			redMatches:    0,
			blueMatches:   1,
			expectedWin:   "八等奖",
			expectedPrize: 500,
			description:   "只选中1个后区",
		},
		// 未中奖情况
		{
			name:          "DLT 未中奖 (1红+0蓝)",
			redMatches:    1,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
			description:   "只选中1个前区",
		},
		{
			name:          "DLT 未中奖 (0红+0蓝)",
			redMatches:    0,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
			description:   "一个都没选中",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winLevel, prizeAmount := calculateWinLevel("dlt", tt.redMatches, tt.blueMatches)
			assert.Equal(t, tt.expectedWin, winLevel, "中奖等级不匹配: %s", tt.description)
			assert.Equal(t, tt.expectedPrize, prizeAmount, "奖金金额不匹配: %s", tt.description)
			t.Logf("✓ %s - %s", tt.name, tt.description)
		})
	}
}

// TestWinningRulesSummary 输出中奖规则汇总
func TestWinningRulesSummary(t *testing.T) {
	t.Log("\n=== 双色球中奖规则汇总 ===")
	t.Log("一等奖: 6红+1蓝 = 500万元")
	t.Log("二等奖: 6红+0蓝 = 10万元")
	t.Log("三等奖: 5红+1蓝 = 3000元")
	t.Log("四等奖: 5红+0蓝 或 4红+1蓝 = 200元")
	t.Log("五等奖: 4红+0蓝 或 3红+1蓝 = 10元")
	t.Log("六等奖: 2红+1蓝 或 1红+1蓝 或 0红+1蓝 = 5元")

	t.Log("\n=== 大乐透中奖规则汇总 ===")
	t.Log("一等奖: 5前+2后 = 1000万元")
	t.Log("二等奖: 5前+1后 = 80万元")
	t.Log("三等奖: 5前+0后 = 1万元")
	t.Log("四等奖: 4前+2后 = 3000元")
	t.Log("五等奖: 4前+1后 或 3前+2后 = 300元")
	t.Log("六等奖: 4前+0后 或 3前+1后 或 2前+2后 = 100元")
	t.Log("七等奖: 3前+0后 或 2前+1后 或 1前+2后 或 0前+2后 = 15元")
	t.Log("八等奖: 2前+0后 或 1前+1后 或 0前+1后 = 5元")
}
