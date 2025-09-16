package fucai

var FucaiHandlerInst IFucai

type IFucai interface {
	GetSSQHistory(req SSQHistoryReq) (res SSQHistoryResp, err error)
}

type FucaiHandler struct {
}

func init() {
	FucaiHandlerInst = NewFucaiHandler()
}

func NewFucaiHandler() *FucaiHandler {
	return &FucaiHandler{}
}
