package http500

// Package http500 提供从500.com网站抓取彩票数据的功能

// GetDLTMissingData 获取大乐透遗漏数据
// periodCount: 期数范围，支持10/30/50
// 返回红球和蓝球的遗漏数据
func GetDLTMissingData(periodCount int) ([]DLTRedMissingData, []DLTBlueMissingData, error) {
	redData, err := FetchDLTRedMissingData(periodCount)
	if err != nil {
		return nil, nil, err
	}

	blueData, err := FetchDLTBlueMissingData(periodCount)
	if err != nil {
		return redData, nil, err
	}

	return redData, blueData, nil
}

// GetSSQMissingData 获取双色球遗漏数据
// periodCount: 期数范围，支持10/30/50
// 返回红球和蓝球的遗漏数据
func GetSSQMissingData(periodCount int) ([]SSQRedMissingData, []SSQBlueMissingData, error) {
	redData, err := FetchSSQRedMissingData(periodCount)
	if err != nil {
		return nil, nil, err
	}

	blueData, err := FetchSSQBlueMissingData(periodCount)
	if err != nil {
		return redData, nil, err
	}

	return redData, blueData, nil
}
