package api

import (
	"net/http"
	"strconv"

	"lucky/common/mysql"
	"lucky/model"
	"lucky/service"

	"github.com/gin-gonic/gin"
)

// CheckWinningRequest 中奖核对请求
type CheckWinningRequest struct {
	NumberID int64 `json:"numberId" binding:"required"` // 用户号码ID
}

// WinningMatch 单次中奖匹配记录
type WinningMatch struct {
	Period      string            `json:"period"`      // 期号
	DrawDate    string            `json:"drawDate"`    // 开奖日期
	RedBalls    model.NumberArray `json:"redBalls"`    // 开奖红球
	BlueBalls   model.NumberArray `json:"blueBalls"`   // 开奖蓝球
	RedMatches  int               `json:"redMatches"`  // 红球匹配数
	BlueMatches int               `json:"blueMatches"` // 蓝球匹配数
	WinLevel    string            `json:"winLevel"`    // 中奖等级
	PrizeAmount int64             `json:"prizeAmount"` // 奖金金额(分)
}

// CheckWinningResponse 中奖核对响应
type CheckWinningResponse struct {
	UserNumber   *model.UserNumber `json:"userNumber"`   // 用户号码信息
	Matches      []WinningMatch    `json:"matches"`      // 中奖匹配列表
	TotalMatches int               `json:"totalMatches"` // 总中奖次数
	TotalPrize   int64             `json:"totalPrize"`   // 总奖金(分)
}

// CheckWinning 核对中奖情况
// @Summary 核对用户号码在近15期的中奖情况
// @Description 对比用户号码和近15期开奖号码，返回中奖的期数、中奖等级等信息
// @Tags 号码管理
// @Accept json
// @Produce json
// @Param numberId path int true "用户号码ID"
// @Success 200 {object} CheckWinningResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/numbers/{numberId}/check [get]
func CheckWinning(c *gin.Context) {
	// 获取用户号码ID
	numberIDStr := c.Param("numberId")
	numberID, err := strconv.ParseInt(numberIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "无效的号码ID",
			"error":   err.Error(),
		})
		return
	}

	// 获取用户号码信息
	userNumberDAO := model.NewUserNumberDAO(mysql.DB)
	userNumber, err := userNumberDAO.GetByIDWithGame(numberID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    404,
			"message": "号码不存在",
			"error":   err.Error(),
		})
		return
	}

	// 获取该游戏的近15期开奖结果
	drawResults, err := service.GetLatestDrawResults(mysql.DB, userNumber.Game.GameCode, 15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取开奖结果失败",
			"error":   err.Error(),
		})
		return
	}

	// 逐一对比中奖情况
	var matches []WinningMatch
	var totalPrize int64

	for _, drawResult := range drawResults {
		redMatches := countMatches(userNumber.RedBalls, drawResult.RedBalls)
		blueMatches := countMatches(userNumber.BlueBalls, drawResult.BlueBalls)

		// 判断是否中奖
		winLevel, prizeAmount := calculateWinLevel(userNumber.Game.GameCode, redMatches, blueMatches)

		if winLevel != "" {
			matches = append(matches, WinningMatch{
				Period:      drawResult.Period,
				DrawDate:    drawResult.DrawDate.Format("2006-01-02"),
				RedBalls:    drawResult.RedBalls,
				BlueBalls:   drawResult.BlueBalls,
				RedMatches:  redMatches,
				BlueMatches: blueMatches,
				WinLevel:    winLevel,
				PrizeAmount: prizeAmount,
			})
			totalPrize += prizeAmount
		}
	}

	response := CheckWinningResponse{
		UserNumber:   userNumber,
		Matches:      matches,
		TotalMatches: len(matches),
		TotalPrize:   totalPrize,
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查询成功",
		"data":    response,
	})
}

// countMatches 计算匹配的号码数量
func countMatches(userBalls, drawBalls model.NumberArray) int {
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
func calculateWinLevel(gameCode string, redMatches, blueMatches int) (string, int64) {
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
