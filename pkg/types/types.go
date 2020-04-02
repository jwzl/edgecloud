package types

import(
	"github.com/jwzl/edgeOn/common"
)

const (
	EDGECLOUD_APISERVER_MODULE	= "apiserver"
	EDGECLOUD_EVENTHUB_MODULE	= "eventhub"
	EDGECLOUD_DEVICETWIN_MODULE = "deviceTwin"
)

type MsgContent struct{
	ReplyChn	chan common.TwinResponse
	Content		interface{}
}

func BuildMessageResponse(code int, reason string, twins []common.DigitalTwin) *common.TwinResponse {

	return &common.TwinResponse{
		Code: code,
		Reason: reason,
		Twins: twins,
	}
}
