package http500

import (
	"encoding/json"
	"testing"
)

func TestIntegration(t *testing.T) {
	t.Run("DLT Missing Data Integration", func(t *testing.T) {
		// 测试大乐透统一接口
		redData, blueData, err := GetDLTMissingData(10)
		if err != nil {
			t.Fatalf("GetDLTMissingData failed: %v", err)
		}
		if len(redData) == 0 {
			t.Fatal("GetDLTMissingData returned empty red data")
		}
		if len(blueData) == 0 {
			t.Fatal("GetDLTMissingData returned empty blue data")
		}
		t.Logf("DLT red balls: %d items", len(redData))
		t.Logf("DLT blue balls: %d items", len(blueData))

		// 验证JSON序列化
		redJSON, _ := json.Marshal(redData[:3])
		blueJSON, _ := json.Marshal(blueData[:3])
		t.Logf("Sample red JSON: %s", string(redJSON))
		t.Logf("Sample blue JSON: %s", string(blueJSON))
	})

	t.Run("SSQ Missing Data Integration", func(t *testing.T) {
		// 测试双色球统一接口
		redData, blueData, err := GetSSQMissingData(10)
		if err != nil {
			t.Fatalf("GetSSQMissingData failed: %v", err)
		}
		if len(redData) == 0 {
			t.Fatal("GetSSQMissingData returned empty red data")
		}
		if len(blueData) == 0 {
			t.Fatal("GetSSQMissingData returned empty blue data")
		}
		t.Logf("SSQ red balls: %d items", len(redData))
		t.Logf("SSQ blue balls: %d items", len(blueData))

		// 验证JSON序列化
		redJSON, _ := json.Marshal(redData[:3])
		blueJSON, _ := json.Marshal(blueData[:3])
		t.Logf("Sample red JSON: %s", string(redJSON))
		t.Logf("Sample blue JSON: %s", string(blueJSON))
	})
}
