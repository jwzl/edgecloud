package apis

import (
	"time"
	"errors"
	"k8s.io/klog"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/edgecloud/pkg/types"
	"github.com/jwzl/beehive/pkg/core/context"
	"github.com/jwzl/edgecloud/pkg/devicetwin/eventlistener"
)

type DeviceTwinModule struct {
	context 				*context.Context
}
var devTM *DeviceTwinModule

func NewDeviceTwinModule(ctx *context.Context){
	devTM = &DeviceTwinModule{
		context: ctx,
	}
}

func (dtm *DeviceTwinModule) SendMessage(edgeID, target, operation, resource string, content interface{}) chan common.TwinResponse {
	resource = edgeID+"/"+resource

	msgContent := types.MsgContent{
		ReplyChn: make(chan common.TwinResponse, 1),
		Content: content,
	}
	modelMsg := common.BuildModelMessage(common.CloudName, target, 
					operation, resource, msgContent)
	/*
	* rewrite the content since the return result is byte array.
	* we just pass the pointer.
	*/
	modelMsg.Content = msgContent	

	dtm.context.Send(types.EDGECLOUD_DEVICETWIN_MODULE, modelMsg)

	return msgContent.ReplyChn
}


// verb: bind
//path: edge/bind?edgeid=xxx 
func BindEdge(edgeID string) (int, string) {
	
	replyChn := devTM.SendMessage(edgeID, common.TwinModuleName, "Bind", common.DGTWINS_RESOURCE_EDGE, nil) 
	resp , ok := <- replyChn
	defer close(replyChn)
	if !ok {
		return common.InternalErrorCode, "Channel has closed"
	}

	//wait the edge online.
	err := eventlistener.WatchEvent(edgeID, "", eventlistener.EVENT_EDGE_ONLINE, 
			500 * time.Millisecond, nil)
	if err != nil {
		return common.InternalErrorCode, err.Error()
	}

	return resp.Code, resp.Reason
}

//path: /edge/twin?edgeid=xxx&twinid=xxx
func CreateTwin(edgeID, twinID string) (int, string){
	replyChn := devTM.SendMessage(edgeID, common.TwinModuleName,
			 common.DGTWINS_OPS_CREATE, common.DGTWINS_RESOURCE_TWINS, twinID)

	resp , ok := <- replyChn
	defer close(replyChn)
	if !ok {
		return common.InternalErrorCode, "Channel has closed"
	}

	//wait the edge online.
	err := eventlistener.WatchEvent(edgeID, twinID, eventlistener.EVENT_TWIN_ONLINE, 
			500 * time.Millisecond, nil)
	
	if err != nil {
		klog.Infof("sdsd %s",err.Error())	
		return common.InternalErrorCode, err.Error()
	}
	return resp.Code, resp.Reason	 
}

//path: /edge/twin?edgeid=xxx&twinid=xxx
func DeleteTwin(edgeID, twinID string) error{
	replyChn := devTM.SendMessage(edgeID, common.TwinModuleName,
			 common.DGTWINS_OPS_DELETE, common.DGTWINS_RESOURCE_TWINS, twinID)
	resp , ok := <- replyChn
	
	if !ok {
		return errors.New("Channel has closed")
	}

	if resp.Code != common.RequestSuccessCode {
		return errors.New(resp.Reason)
	}

	return nil
}

//path: /edge/twin?edgeid=xxx&twinid=xxx
func UpdateProperty(edgeID, twinID string, desired map[string]*common.TwinProperty) error {

	dgtwin := &common.DigitalTwin{
		ID: twinID,	
		Properties: common.TwinProperties{
			Desired: desired,
		},
	}
	twins := []common.DigitalTwin{*dgtwin}

	msgContent, err := common.BuildTwinMessage(twins)
	if err != nil {
		return err
	}
	replyChn := devTM.SendMessage(edgeID, common.TwinModuleName,
			 common.DGTWINS_OPS_GET, common.DGTWINS_RESOURCE_TWINS, msgContent)
	resp , ok := <- replyChn
	if !ok {
		return errors.New("Channel has closed")
	}

	if resp.Code != common.RequestSuccessCode {
		return errors.New(resp.Reason)
	}

	return nil
}


//path: /edge/twin?edgeid=xxx&twinid=xxx
func GetTwin(edgeID, twinID string) (*common.DigitalTwin, error){
	replyChn := devTM.SendMessage(edgeID, common.TwinModuleName,
			 common.DGTWINS_OPS_GET, common.DGTWINS_RESOURCE_TWINS, []string{twinID})

	resp , ok := <- replyChn
	if !ok {
		return nil, errors.New("Channel has closed")
	}

	if resp.Code != common.RequestSuccessCode {
		return nil, errors.New(resp.Reason)
	}

	return &resp.Twins[0], nil
} 


