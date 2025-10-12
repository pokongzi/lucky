package http500

import (
	"fmt"
	"testing"
)

func TestFetchDLTRedMissingData(t *testing.T) {
	tests := []struct {
		name        string
		periodCount int
	}{
		{
			name:        "期数10",
			periodCount: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := FetchDLTRedMissingData(tt.periodCount)
			if err != nil {
				t.Errorf("FetchDLTRedMissingData(%d) error = %v", tt.periodCount, err)
				return
			}

			// 检查数据是否有效
			if len(data) == 0 {
				t.Errorf("FetchDLTRedMissingData(%d) returned empty data", tt.periodCount)
			} else {
				t.Logf("FetchDLTRedMissingData(%d) returned %d items", tt.periodCount, len(data))

				// 输出部分数据，用于检查
				for i := 0; i < 3 && i < len(data); i++ {
					t.Logf("Sample data %d: %+v", i, data[i])
				}
			}

			// 测试JSON序列化
			json, err := GetDLTRedMissingDataJSON(tt.periodCount)
			if err != nil {
				t.Errorf("GetDLTRedMissingDataJSON(%d) error = %v", tt.periodCount, err)
				return
			}

			t.Logf("JSON sample: %s", json[:100]+"...")
		})
	}
}

func TestFetchDLTBlueMissingData(t *testing.T) {
	tests := []struct {
		name        string
		periodCount int
	}{
		{
			name:        "期数10",
			periodCount: 10,
		},
		{
			name:        "期数30",
			periodCount: 30,
		},
		{
			name:        "期数50",
			periodCount: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := FetchDLTBlueMissingData(tt.periodCount)
			if err != nil {
				t.Errorf("FetchDLTBlueMissingData(%d) error = %v", tt.periodCount, err)
				return
			}

			// 检查数据是否有效
			if len(data) == 0 {
				t.Errorf("FetchDLTBlueMissingData(%d) returned empty data", tt.periodCount)
			} else {
				t.Logf("FetchDLTBlueMissingData(%d) returned %d items", tt.periodCount, len(data))

				// 输出部分数据，用于检查
				for i := 0; i < 3 && i < len(data); i++ {
					t.Logf("Sample data %d: %+v", i, data[i])
				}
			}

			// 测试JSON序列化
			json, err := GetDLTBlueMissingDataJSON(tt.periodCount)
			if err != nil {
				t.Errorf("GetDLTBlueMissingDataJSON(%d) error = %v", tt.periodCount, err)
				return
			}

			t.Logf("JSON sample: %s", json[:100]+"...")
		})
	}
}

// Example 提供了使用示例
func Example() {
	// 获取大乐透红球30期遗漏数据
	redData, err := FetchDLTRedMissingData(30)
	if err != nil {
		fmt.Printf("获取大乐透红球遗漏数据失败: %v\n", err)
		return
	}

	fmt.Printf("大乐透红球遗漏数据共 %d 条\n", len(redData))

	// 获取大乐透蓝球30期遗漏数据
	blueData, err := FetchDLTBlueMissingData(30)
	if err != nil {
		fmt.Printf("获取大乐透蓝球遗漏数据失败: %v\n", err)
		return
	}

	fmt.Printf("大乐透蓝球遗漏数据共 %d 条\n", len(blueData))

	// 获取双色球红球30期遗漏数据
	ssqRedData, err := FetchSSQRedMissingData(30)
	if err != nil {
		fmt.Printf("获取双色球红球遗漏数据失败: %v\n", err)
		return
	}

	fmt.Printf("双色球红球遗漏数据共 %d 条\n", len(ssqRedData))

	// 获取双色球蓝球30期遗漏数据
	ssqBlueData, err := FetchSSQBlueMissingData(30)
	if err != nil {
		fmt.Printf("获取双色球蓝球遗漏数据失败: %v\n", err)
		return
	}

	fmt.Printf("双色球蓝球遗漏数据共 %d 条\n", len(ssqBlueData))
}

func TestFetchSSQRedMissingData(t *testing.T) {
	tests := []struct {
		name        string
		periodCount int
	}{
		{
			name:        "期数10",
			periodCount: 10,
		},
		{
			name:        "期数30",
			periodCount: 30,
		},
		{
			name:        "期数50",
			periodCount: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := FetchSSQRedMissingData(tt.periodCount)
			if err != nil {
				t.Errorf("FetchSSQRedMissingData(%d) error = %v", tt.periodCount, err)
				return
			}

			// 检查数据是否有效
			if len(data) == 0 {
				t.Errorf("FetchSSQRedMissingData(%d) returned empty data", tt.periodCount)
			} else {
				t.Logf("FetchSSQRedMissingData(%d) returned %d items", tt.periodCount, len(data))

				// 输出部分数据，用于检查
				for i := 0; i < 3 && i < len(data); i++ {
					t.Logf("Sample data %d: %+v", i, data[i])
				}
			}

			// 测试JSON序列化
			json, err := GetSSQRedMissingDataJSON(tt.periodCount)
			if err != nil {
				t.Errorf("GetSSQRedMissingDataJSON(%d) error = %v", tt.periodCount, err)
				return
			}

			t.Logf("JSON sample: %s", json[:100]+"...")
		})
	}
}

func TestFetchSSQBlueMissingData(t *testing.T) {
	tests := []struct {
		name        string
		periodCount int
	}{
		{
			name:        "期数10",
			periodCount: 10,
		},
		{
			name:        "期数30",
			periodCount: 30,
		},
		{
			name:        "期数50",
			periodCount: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := FetchSSQBlueMissingData(tt.periodCount)
			if err != nil {
				t.Errorf("FetchSSQBlueMissingData(%d) error = %v", tt.periodCount, err)
				return
			}

			// 检查数据是否有效
			if len(data) == 0 {
				t.Errorf("FetchSSQBlueMissingData(%d) returned empty data", tt.periodCount)
			} else {
				t.Logf("FetchSSQBlueMissingData(%d) returned %d items", tt.periodCount, len(data))

				// 输出部分数据，用于检查
				for i := 0; i < 3 && i < len(data); i++ {
					t.Logf("Sample data %d: %+v", i, data[i])
				}
			}

			// 测试JSON序列化
			json, err := GetSSQBlueMissingDataJSON(tt.periodCount)
			if err != nil {
				t.Errorf("GetSSQBlueMissingDataJSON(%d) error = %v", tt.periodCount, err)
				return
			}

			t.Logf("JSON sample: %s", json[:100]+"...")
		})
	}
}
