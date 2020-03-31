package apis

import (
	"errors"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/edgecloud/pkg/types"
	"github.com/jwzl/beehive/pkg/core/context"	
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

	dtm.context.Send(types.EDGECLOUD_DEVICETWIN_MODULE, modelMsg)

	return msgContent.ReplyChn
}


// verb: bind
//path: edge/bind?edgeid=xxx 
func BindEdge(edgeID string) error {
	
	replyChn := devTM.SendMessage(edgeID, "edge", "Bind", common.DGTWINS_RESOURCE_EDGE, nil) 
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
func CreateTwin(edgeID, twinID string) error{
	replyChn := devTM.SendMessage(edgeID, common.TwinModuleName,
			 common.DGTWINS_OPS_CREATE, common.DGTWINS_RESOURCE_TWINS, twinID)

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
			 common.DGTWINS_OPS_GET, common.DGTWINS_RESOURCE_TWINS, twinID)

	resp , ok := <- replyChn
	if !ok {
		return nil, errors.New("Channel has closed")
	}

	if resp.Code != common.RequestSuccessCode {
		return nil, errors.New(resp.Reason)
	}

	return &resp.Twins[0], nil
} 


