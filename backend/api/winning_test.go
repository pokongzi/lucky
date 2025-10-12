package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCheckWinningValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建路由
	r := gin.New()
	r.GET("/api/numbers/:numberId/check", CheckWinning)

	tests := []struct {
		name         string
		numberId     string
		expectedCode int
		errorMessage string
	}{
		{
			name:         "Invalid number ID - not numeric",
			numberId:     "abc",
			expectedCode: http.StatusBadRequest,
			errorMessage: "无效的号码ID",
		},
		{
			name:         "Invalid number ID - negative",
			numberId:     "-1",
			expectedCode: http.StatusBadRequest,
			errorMessage: "无效的号码ID",
		},
		{
			name:         "Invalid number ID - zero",
			numberId:     "0",
			expectedCode: http.StatusNotFound, // 因为数据库中不存在ID为0的记录
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			req, _ := http.NewRequest("GET", "/api/numbers/"+tt.numberId+"/check", nil)
			req.Header.Set("X-User-ID", "1")

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证状态码
			assert.Equal(t, tt.expectedCode, w.Code)

			// 验证错误消息（如果有）
			if tt.errorMessage != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["message"], tt.errorMessage)
			}
		})
	}
}

func TestCountMatches(t *testing.T) {
	tests := []struct {
		name      string
		userBalls []int
		drawBalls []int
		expected  int
	}{
		{
			name:      "No matches",
			userBalls: []int{1, 2, 3, 4, 5},
			drawBalls: []int{6, 7, 8, 9, 10},
			expected:  0,
		},
		{
			name:      "Partial matches",
			userBalls: []int{1, 2, 3, 4, 5},
			drawBalls: []int{1, 2, 8, 9, 10},
			expected:  2,
		},
		{
			name:      "All matches",
			userBalls: []int{1, 2, 3, 4, 5},
			drawBalls: []int{1, 2, 3, 4, 5},
			expected:  5,
		},
		{
			name:      "Duplicate in user balls",
			userBalls: []int{1, 1, 2, 3, 4},
			drawBalls: []int{1, 5, 6, 7, 8},
			expected:  2, // 两个1都匹配同一个开奖号码1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countMatches(tt.userBalls, tt.drawBalls)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalculateWinLevel(t *testing.T) {
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
			name:          "SSQ Second Prize",
			gameCode:      "ssq",
			redMatches:    6,
			blueMatches:   0,
			expectedWin:   "二等奖",
			expectedPrize: 10000000,
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
			name:          "DLT Eighth Prize",
			gameCode:      "dlt",
			redMatches:    0,
			blueMatches:   1,
			expectedWin:   "八等奖",
			expectedPrize: 500,
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
			winLevel, prizeAmount := calculateWinLevel(tt.gameCode, tt.redMatches, tt.blueMatches)
			assert.Equal(t, tt.expectedWin, winLevel)
			assert.Equal(t, tt.expectedPrize, prizeAmount)
		})
	}
}
