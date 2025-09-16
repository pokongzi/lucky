package fucai

// SSQHistoryReq 双色球历史数据请求参数
type SSQHistoryReq struct {
	Name       string `json:"name"`       // 游戏名称，双色球为ssq
	IssueCount string `json:"issueCount"` // 期数
	IssueStart string `json:"issueStart"` // 开始期号
	IssueEnd   string `json:"issueEnd"`   // 结束期号
	DayStart   string `json:"dayStart"`   // 开始日期
	DayEnd     string `json:"dayEnd"`     // 结束日期
	PageNo     int    `json:"pageNo"`     // 页码
	PageSize   int    `json:"pageSize"`   // 每页数据量
	Week       string `json:"week"`       // 周
	SystemType string `json:"systemType"` // 系统类型，PC
}

// SSQHistoryItem 单期双色球数据
type SSQHistoryItem struct {
	Code string `json:"code"` // 期号
	Date string `json:"date"` // 开奖日期
	Red  string `json:"red"`  // 红球，格式：01,05,16,20,21,32
	Blue string `json:"blue"` // 蓝球
}

// SSQHistoryResp 双色球历史数据响应
type SSQHistoryResp struct {
	State   int              `json:"state"`
	Message string           `json:"message"`
	Total   int              `json:"total"`
	TFoot   string           `json:"Tfooter"`
	Result  []SSQHistoryItem `json:"result"`
}
