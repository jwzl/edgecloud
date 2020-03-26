package devicetwin

import(
	"time"	
	"strings"
	"encoding/json"

	"k8s.io/klog"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/jwzl/edgecloud/pkg/types"
	"github.com/jwzl/beehive/pkg/core/context"
)

const (
	DGTWINS_EDGE_BIND	= "Bind"

	MODEL_MSG_TIMEOUT = 5*60		//5s 
)

type DeviceTwinModule struct {
	DownStreamMessageChan	chan *model.Message
	context 				*context.Context
	dtcontext 				*DTContext
}

// Register this module.
func Register(){	
	dtm := &DeviceTwinModule{}
	core.Register(dtm)
}

//Name
func (dtm *DeviceTwinModule) Name() string {
	return "deviceTwin"
}

//Group
func (dtm *DeviceTwinModule) Group() string {
	return "deviceTwin"
}

//Start this module.
func (dtm *DeviceTwinModule) Start(c *context.Context) {

	dtm.context = c
	dtm.dtcontext = NewDTContext(c)	

	go dtm.MessageCheck()

	for {
		v, err := dtm.context.Receive(dtm.Name())
		if err != nil {
			klog.Errorf("failed to receive message: %v", err)
			break
		}

		msg, isThisType := v.(*model.Message)
		if !isThisType || msg == nil { 		
			continue
		}

		target := msg.GetTarget()
		operation := msg.GetOperation()
		
		if strings.Contains(operation, common.DGTWINS_OPS_KEEPALIVE){
			// recieve the heartbeat.
			dtm.dtcontext.DealHeartBeat(msg)
		}else if strings.Contains(target, common.CloudName) {
			// Do Up stream message.
			go dtm.doUpStreamMessage(msg)
		}else if strings.Contains(target, common.TwinModuleName) {
			// Do Down stream message.
			go dtm.doDownStreamMessage(msg)
		}
	}
	
	
}

//Cleanup
func (dtm *DeviceTwinModule) Cleanup() {
	dtm.context.Cleanup(dtm.Name())
}

/*
* doUpStreamMessage
*/
func (dtm *DeviceTwinModule) doUpStreamMessage(msg *model.Message) {
	operation := msg.GetOperation()
	resource  := msg.GetResource()
	splitString := strings.Split(resource, "/")	
	edgeID := splitString[0]

	switch operation {
	case common.DGTWINS_OPS_RESPONSE:
		if strings.Contains(splitString[1], common.DGTWINS_RESOURCE_EDGE) {
			/*
			* This is a edge information report.
			*/
			contents, isThisType := msg.GetContent().([]byte)
			if !isThisType {
				return 
			}
		
			var edgeInfo common.EdgeInfo
			err := json.Unmarshal(contents, &edgeInfo)
			if err != nil {
				return 
			}

			//build a edge description struct.
			edged:= NewEdgeDescription(edgeInfo.EdgeID)
			edged.SetEdgeName(edgeInfo.EdgeName)
			edged.SetEdgeDescription(edgeInfo.Description)
			edged.SetEdgeState(EdgeStateOnline) 

			err = dtm.dtcontext.AddEdgeInfo(edged)
			if err != nil {
				klog.Infof("Add edge info failed")
			}
		}else{
			//If no cache such message, then, Ignore this message. 
			if dtm.dtcontext.CacheHasThisMessage(msg) != true {
				return
			}	

			respMsg, err := common.UnMarshalResponseMessage(msg)
			if err != nil {
				return
			}	

			if respMsg.Code == common.RequestSuccessCode {
				// delete this message.
				dtm.dtcontext.DeleteMsgCache(msg)
			}else {
				klog.Warningf("Unexpected message:", respMsg)	
			}
		}
	case common.DGTWINS_OPS_SYNC:
		twinMsg, err := common.UnMarshalTwinMessage(msg)
		if err != nil {
			return
		}	

		dgTwin := &twinMsg.Twins[0]
		//build response message.
		twin := &common.DigitalTwin{
			ID: dgTwin.ID,
		}
		twins := []common.DigitalTwin{*twin}

		err = dtm.dtcontext.UpdateTwin(edgeID, dgTwin)
		if err != nil {
			return
		}	
		klog.Infof("Update successful.") 
		
		msgContent, err := common.BuildResponseMessage(common.RequestSuccessCode, "Update successful", twins)
		if err != nil {
			return 
		}
		dtm.dtcontext.SendResponseMessage(msg, msgContent)
	}	
}

/*
* doDownStreamMessage
*/
func (dtm *DeviceTwinModule) doDownStreamMessage(msg *model.Message) {
	//var err error

	operation := msg.GetOperation()
	resource := msg.GetResource()
	splitString := strings.Split(resource, "/")
	edgeID := splitString[0]
	msgContent, isThisType := msg.GetContent().(types.MsgContent)
	if !isThisType {
		//ignore.
		return 
	}
	replyChn := msgContent.ReplyChn	

	switch operation {
	case DGTWINS_EDGE_BIND:
		msg.Router.Target = "edge"
		//send to event hub.
		dtm.context.Send("EventHub", msg)
		//cache the message.
		dtm.dtcontext.CacheMessage(msg)	
		
		//reply the message.
		resp := types.BuildMessageResponse(common.RequestSuccessCode, "successful", nil)
		replyChn <- *resp
	case common.DGTWINS_OPS_CREATE:
		twinID, isthisType := msg.GetContent().(string)
		if !isthisType {
			klog.Warningf("Error format")
			return
		}
		err := dtm.dtcontext.RegisterTwins(edgeID, twinID)
		if err != nil {
			klog.Warningf("register failed")
			return
		}
		
		// response the successful message.
		//reply the message.
		resp := types.BuildMessageResponse(common.RequestSuccessCode, "Update successful", nil)
		replyChn <- *resp

		// Send create twin message.
		twin := &common.DigitalTwin{
			ID: twinID,
		}
		twins := []common.DigitalTwin{*twin}
	 	msgContent, err := common.BuildTwinMessage(twins)
		if err != nil {
			return 
		}

		dtm.dtcontext.SendTwinMessage(edgeID, common.DGTWINS_OPS_CREATE, msgContent)
		//cache the message.
		dtm.dtcontext.CacheMessage(msg)	 
	case common.DGTWINS_OPS_UPDATE:	
		/*
		* update the twin desired property.
		*/
		twinMsg, err := common.UnMarshalTwinMessage(msg)
		if err != nil {
			return
		}	
		dgTwin := &twinMsg.Twins[0]
		//update the local twin property.
		err = dtm.dtcontext.UpdateTwin(edgeID, dgTwin)
		if err != nil {
			return
		}	
		klog.Infof("Update successful.") 
		dtm.dtcontext.SendPropertyMessage(edgeID, common.DGTWINS_OPS_UPDATE, msg.GetContent())
		//cache the message.
		dtm.dtcontext.CacheMessage(msg)	
		//reply the message.
		resp := types.BuildMessageResponse(common.RequestSuccessCode, "Update successful", nil)
		replyChn <- *resp
	case common.DGTWINS_OPS_DELETE:
		/*
		* delete the twin by id.
		*/
		twinID, isthisType := msg.GetContent().(string)
		if !isthisType {
			klog.Warningf("Error format")
			return
		}
		err := dtm.dtcontext.DeleteTwin(edgeID, twinID)
		if err != nil {
			klog.Warningf("register failed")
			return
		}

		//build the send message.
		twin := &common.DigitalTwin{
			ID: twinID,
		}
		twins := []common.DigitalTwin{*twin}
	 	msgContent, err := common.BuildTwinMessage(twins)
		if err != nil {
			return 
		}

		dtm.dtcontext.SendTwinMessage(edgeID, common.DGTWINS_OPS_DELETE, msgContent)
		//cache the message.
		dtm.dtcontext.CacheMessage(msg)	
		//reply the message.
		resp := types.BuildMessageResponse(common.RequestSuccessCode, "Update successful", nil)
		replyChn <- *resp			
	case common.DGTWINS_OPS_GET:
		/*
		* Get the current twin.
		*/
		twinIDs, isthisType := msg.GetContent().([]string)
		if !isthisType {
			klog.Warningf("Error format")
			return
		}

		twins := make([]common.DigitalTwin, 0)
		for _, twinID := range twinIDs {
			content, err := dtm.dtcontext.GetRawTwin(edgeID, twinID)
			if err != nil {
				klog.Warningf("Get failed")
				continue
			} 

			var dgTwin common.DigitalTwin
			err = json.Unmarshal(content, &dgTwin)
			if err != nil {
				return 
			}
			twins = append(twins, dgTwin)
		}

		//reply the message.
		resp := types.BuildMessageResponse(common.RequestSuccessCode, "Get", twins)
		replyChn <- *resp
	default:
		klog.Warningf("Ignored message:", msg)
		resp := types.BuildMessageResponse(common.BadRequestCode, "Ignored", nil)
		replyChn <- *resp
	}
}

func (dtm *DeviceTwinModule) MessageCheck(){
	checkTimeoutCh := time.After(10*time.Second)
	checkHealthCh  := time.After(60*time.Second)
	for {
		select {
		case <-checkHealthCh:
			//health check.
			dtm.dtcontext.EdgeHealth.Range(func (key interface{}, value interface{}) bool {
				edgeID := key.(string)
				timeStamp := value.(int64)
				now := time.Now().Unix()
				if now - timeStamp > 80 {
					klog.Infof("edge %s is not healthy, we mark the edge as offline", edgeID)
					dtm.dtcontext.SetEdgeState(edgeID, EdgeStateOffline)
				}else{
					dtm.dtcontext.SetEdgeState(edgeID, EdgeStateOnline)
				}

				return true
			})

			checkHealthCh = time.After(60*time.Second)
		case <-checkTimeoutCh:
			//check  the MessageCache for response.
			dtm.dealMessageTimeout()	
			checkTimeoutCh = time.After(10*time.Second)
		}
	}		
}

func (dtm *DeviceTwinModule) dealMessageTimeout() {
	dtm.dtcontext.MessageCache.Range(func (key interface{}, value interface{}) bool {
		msg, isMsgType := value.(*model.Message)
		if !isMsgType {
			return false
		}

		timeStamp := msg.GetTimestamp()/1e3
		now	:= time.Now().UnixNano() / 1e9
		if now - timeStamp >= MODEL_MSG_TIMEOUT {
			/*
			* Timeout, delete the message. 
			*/
			dtm.dtcontext.MessageCache.Delete(key)
			return true
		}else{
			//resend  the message.
			klog.Infof("resend the message:", msg)	
			dtm.context.Send("EventHub", msg)
			return true
		} 

		return false
	})
}
