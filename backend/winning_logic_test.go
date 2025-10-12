package main

import (
	"fmt"
	"testing"
)

// NumberArray 模拟类型
type NumberArray []int

// countMatches 计算匹配的号码数量
func countMatchesTest(userBalls, drawBalls NumberArray) int {
	matches := 0
	ballMap := make(map[int]bool)

	// 将开奖号码放入map中
	for _, ball := range drawBalls {
		ballMap[ball] = true
	}

	// 检查用户号码中有多少匹配
	for _, ball := range userBalls {
		if ballMap[ball] {
			matches++
		}
	}

	return matches
}

// calculateWinLevel 根据匹配数计算中奖等级和奖金
func calculateWinLevelTest(gameCode string, redMatches, blueMatches int) (string, int64) {
	if gameCode == "ssq" {
		// 双色球中奖规则
		switch {
		case redMatches == 6 && blueMatches == 1:
			return "一等奖", 500000000 // 5百万(模拟金额)
		case redMatches == 6 && blueMatches == 0:
			return "二等奖", 10000000 // 10万
		case redMatches == 5 && blueMatches == 1:
			return "三等奖", 300000 // 3千
		case (redMatches == 5 && blueMatches == 0) || (redMatches == 4 && blueMatches == 1):
			return "四等奖", 20000 // 200元
		case (redMatches == 4 && blueMatches == 0) || (redMatches == 3 && blueMatches == 1):
			return "五等奖", 1000 // 10元
		case (redMatches == 2 && blueMatches == 1) || (redMatches == 1 && blueMatches == 1) || (redMatches == 0 && blueMatches == 1):
			return "六等奖", 500 // 5元
		}
	} else if gameCode == "dlt" {
		// 大乐透中奖规则
		switch {
		case redMatches == 5 && blueMatches == 2:
			return "一等奖", 1000000000 // 1千万(模拟金额)
		case redMatches == 5 && blueMatches == 1:
			return "二等奖", 80000000 // 80万
		case redMatches == 5 && blueMatches == 0:
			return "三等奖", 1000000 // 1万
		case redMatches == 4 && blueMatches == 2:
			return "四等奖", 300000 // 3千
		case (redMatches == 4 && blueMatches == 1) || (redMatches == 3 && blueMatches == 2):
			return "五等奖", 30000 // 300元
		case (redMatches == 4 && blueMatches == 0) || (redMatches == 3 && blueMatches == 1) || (redMatches == 2 && blueMatches == 2):
			return "六等奖", 10000 // 100元
		case (redMatches == 3 && blueMatches == 0) || (redMatches == 2 && blueMatches == 1) || (redMatches == 1 && blueMatches == 2) || (redMatches == 0 && blueMatches == 2):
			return "七等奖", 1500 // 15元
		case (redMatches == 2 && blueMatches == 0) || (redMatches == 1 && blueMatches == 1) || (redMatches == 0 && blueMatches == 1):
			return "八等奖", 500 // 5元
		}
	}

	return "", 0 // 未中奖
}

func TestCountMatchesStandalone(t *testing.T) {
	tests := []struct {
		name      string
		userBalls NumberArray
		drawBalls NumberArray
		expected  int
	}{
		{
			name:      "No matches",
			userBalls: NumberArray{1, 2, 3, 4, 5},
			drawBalls: NumberArray{6, 7, 8, 9, 10},
			expected:  0,
		},
		{
			name:      "Partial matches",
			userBalls: NumberArray{1, 2, 3, 4, 5},
			drawBalls: NumberArray{1, 2, 8, 9, 10},
			expected:  2,
		},
		{
			name:      "All matches",
			userBalls: NumberArray{1, 2, 3, 4, 5},
			drawBalls: NumberArray{1, 2, 3, 4, 5},
			expected:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countMatchesTest(tt.userBalls, tt.drawBalls)
			if result != tt.expected {
				t.Errorf("countMatches() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateWinLevelStandalone(t *testing.T) {
	tests := []struct {
		name          string
		gameCode      string
		redMatches    int
		blueMatches   int
		expectedWin   string
		expectedPrize int64
	}{
		// 双色球测试用例
		{
			name:          "SSQ First Prize",
			gameCode:      "ssq",
			redMatches:    6,
			blueMatches:   1,
			expectedWin:   "一等奖",
			expectedPrize: 500000000,
		},
		{
			name:          "SSQ Sixth Prize",
			gameCode:      "ssq",
			redMatches:    0,
			blueMatches:   1,
			expectedWin:   "六等奖",
			expectedPrize: 500,
		},
		{
			name:          "SSQ No Prize",
			gameCode:      "ssq",
			redMatches:    0,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
		},
		// 大乐透测试用例
		{
			name:          "DLT First Prize",
			gameCode:      "dlt",
			redMatches:    5,
			blueMatches:   2,
			expectedWin:   "一等奖",
			expectedPrize: 1000000000,
		},
		{
			name:          "DLT No Prize",
			gameCode:      "dlt",
			redMatches:    0,
			blueMatches:   0,
			expectedWin:   "",
			expectedPrize: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			winLevel, prizeAmount := calculateWinLevelTest(tt.gameCode, tt.redMatches, tt.blueMatches)
			if winLevel != tt.expectedWin {
				t.Errorf("calculateWinLevel() winLevel = %v, want %v", winLevel, tt.expectedWin)
			}
			if prizeAmount != tt.expectedPrize {
				t.Errorf("calculateWinLevel() prizeAmount = %v, want %v", prizeAmount, tt.expectedPrize)
			}
		})
	}
}

// 手动验证示例
func TestManualExample(t *testing.T) {
	fmt.Println("=== 中奖核对逻辑测试 ===")

	// 测试号码匹配
	userRed := NumberArray{1, 2, 3, 4, 5, 6}
	userBlue := NumberArray{7}
	drawRed := NumberArray{1, 2, 3, 15, 20, 25}
	drawBlue := NumberArray{7}

	redMatches := countMatchesTest(userRed, drawRed)
	blueMatches := countMatchesTest(userBlue, drawBlue)

	fmt.Printf("用户红球: %v\n", userRed)
	fmt.Printf("开奖红球: %v\n", drawRed)
	fmt.Printf("红球匹配: %d个\n", redMatches)
	fmt.Printf("用户蓝球: %v\n", userBlue)
	fmt.Printf("开奖蓝球: %v\n", drawBlue)
	fmt.Printf("蓝球匹配: %d个\n", blueMatches)

	// 测试双色球中奖等级
	winLevel, prizeAmount := calculateWinLevelTest("ssq", redMatches, blueMatches)
	fmt.Printf("双色球中奖等级: %s\n", winLevel)
	fmt.Printf("奖金: %.2f元\n", float64(prizeAmount)/100)

	// 测试大乐透中奖等级
	winLevel2, prizeAmount2 := calculateWinLevelTest("dlt", redMatches, blueMatches)
	fmt.Printf("大乐透中奖等级: %s\n", winLevel2)
	fmt.Printf("奖金: %.2f元\n", float64(prizeAmount2)/100)
}
