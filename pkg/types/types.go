package types

import(
	"github.com/jwzl/edgeOn/common"
)

type MsgContent struct{
	ReplyChn	chan common.TwinResponse
	Content		interface{}
}

func BuildMessageResponse(code int, reason string, twins []DigitalTwin) *common.TwinResponse {

	return &common.TwinResponse{
		Code: code,
		Reason: reason,
		Twins: twins,
	}
}
