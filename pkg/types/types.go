package types

const (
	EDGECLOUD_APISERVER_MODULE	= "apiserver"
	EDGECLOUD_EVENTHUB_MODULE	= "eventhub"
	EDGECLOUD_DEVICETWIN_MODULE = "deviceTwin"
)

/*
* general response structure.
*/
type Response struct{
	Code 	int 			`json:"code"`
	Reason 	string 			`json:"reason,omitempty"`
	Content  interface{}	`json:"twins,omitempty"`	
}

type MsgContent struct{
	ReplyChn	chan Response
	Content		interface{}
}

type EdgeInfo struct{
	ID				string
	Name 			string
	Description		string
	State			string
	/*
	* all device ID in this edge.
	*/
	DeviceIDs  		[]string
}

func BuildResponse(code int, reason string, content interface{}) *Response {
	return &Response{
		Code: code,
		Reason: reason,
		Content: content,
	}
}

/*func BuildMessageResponse(code int, reason string, twins []common.DigitalTwin) *common.TwinResponse {

	return &common.TwinResponse{
		Code: code,
		Reason: reason,
		Twins: twins,
	}
}*/
