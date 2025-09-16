package ticai

var TicaiHandlerInst ITicai

type ITicai interface {
	GetDLTHistory(req DLTHistoryReq) (res DLTHistoryResp, err error)
}

type TicaiHandler struct {
}

func init() {
	TicaiHandlerInst = NewTicaiHandler()
}

func NewTicaiHandler() *TicaiHandler {
	return &TicaiHandler{}
}
