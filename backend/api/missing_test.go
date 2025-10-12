package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMissingDataValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建路由
	r := gin.New()
	r.GET("/api/missing", GetMissingData)

	tests := []struct {
		name         string
		gameCode     string
		periodCount  string
		expectedCode int
		errorMessage string
	}{
		{
			name:         "Invalid game code",
			gameCode:     "invalid",
			periodCount:  "10",
			expectedCode: http.StatusBadRequest,
			errorMessage: "gameCode 只支持 dlt(大乐透) 或 ssq(双色球)",
		},
		{
			name:         "Invalid period count",
			gameCode:     "dlt",
			periodCount:  "20",
			expectedCode: http.StatusBadRequest,
			errorMessage: "periodCount 只支持 10、30、50",
		},
		{
			name:         "Missing game code",
			gameCode:     "",
			periodCount:  "10",
			expectedCode: http.StatusBadRequest,
			errorMessage: "gameCode 参数不能为空",
		},
		{
			name:         "Missing period count",
			gameCode:     "dlt",
			periodCount:  "",
			expectedCode: http.StatusBadRequest,
			errorMessage: "periodCount 必须是数字",
		},
		{
			name:         "Non-numeric period count",
			gameCode:     "dlt",
			periodCount:  "abc",
			expectedCode: http.StatusBadRequest,
			errorMessage: "periodCount 必须是数字",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			req, _ := http.NewRequest("GET", "/api/missing", nil)
			q := req.URL.Query()
			if tt.gameCode != "" {
				q.Add("gameCode", tt.gameCode)
			}
			if tt.periodCount != "" {
				q.Add("periodCount", tt.periodCount)
			}
			req.URL.RawQuery = q.Encode()

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证状态码
			assert.Equal(t, tt.expectedCode, w.Code)

			// 验证错误消息
			if tt.errorMessage != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.errorMessage)
			}
		})
	}
}

func TestGetMissingDataBatchValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// 创建路由
	r := gin.New()
	r.GET("/api/missing/batch", GetMissingDataBatch)

	tests := []struct {
		name         string
		gameCode     string
		expectedCode int
		errorMessage string
	}{
		{
			name:         "Invalid game code",
			gameCode:     "invalid",
			expectedCode: http.StatusBadRequest,
			errorMessage: "gameCode 只支持 dlt(大乐透) 或 ssq(双色球)",
		},
		{
			name:         "Missing game code",
			gameCode:     "",
			expectedCode: http.StatusBadRequest,
			errorMessage: "gameCode 参数不能为空",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建请求
			req, _ := http.NewRequest("GET", "/api/missing/batch", nil)
			if tt.gameCode != "" {
				q := req.URL.Query()
				q.Add("gameCode", tt.gameCode)
				req.URL.RawQuery = q.Encode()
			}

			// 创建响应记录器
			w := httptest.NewRecorder()

			// 执行请求
			r.ServeHTTP(w, req)

			// 验证状态码
			assert.Equal(t, tt.expectedCode, w.Code)

			// 验证错误消息
			if tt.errorMessage != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"], tt.errorMessage)
			}
		})
	}
}
