package ticai

// DLTHistoryReq 大乐透历史数据请求参数
type DLTHistoryReq struct {
	GameNo     string `json:"gameNo"`     // 游戏编号，大乐透为85
	ProvinceId string `json:"provinceId"` // 省份ID，0表示全国
	PageSize   int    `json:"pageSize"`   // 每页数据量
	PageNo     int    `json:"pageNo"`     // 页码
	IsVerify   int    `json:"isVerify"`   // 是否验证，1表示是
}

// DLTHistoryItem 单期大乐透数据
type DLTHistoryItem struct {
	LotteryDrawNum    string `json:"lotteryDrawNum"`    // 期号
	LotteryDrawTime   string `json:"lotteryDrawTime"`   // 开奖日期
	LotteryDrawResult string `json:"lotteryDrawResult"` // 开奖结果，格式：01,11,14,25,27+04,10
}

// DLTHistoryResp 大乐透历史数据响应
type DLTHistoryResp struct {
	Value struct {
		List     []DLTHistoryItem `json:"list"`
		PageNo   int              `json:"pageNo"`
		PageSize int              `json:"pageSize"`
		Pages    int              `json:"pages"`
		Total    int              `json:"total"`
	} `json:"value"`
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
