package devicetwin

const (
	DGTWINS_EDGE_BIND	= "Bind"
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
			edgeInfo, isThisType := msg.GetContent().(common.EdgeInfo)
			if !isThisType {
				return 
			}
	
			//build a edge description struct.
			edged:= NewEdgeDescription(edgeInfo.edgeID)
			edged.SetEdgeName(edgeInfo.EdgeName)
			edged.SetEdgeDescription(edgeInfo.Description)
			edged.SetEdgeState(EdgeStateOnline) 

			err := dtm.dtcontext.AddEdgeInfo(edged)
			if err != nil {
				klog.Infof("Add edge info failed")
			}
		}else{	
			respMsg, err := common.UnMarshalResponseMessage(msg)
			if err != nil {
				return
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

		err := dtm.dtcontext.UpdateTwin(edgeID, dgTwin)
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
	var err error

	operation := msg.GetOperation()
	resource := msg.GetResource()
	splitString := strings.Split(resource, "/")
	edgeID := splitString[0]	

	switch operation {
	case DGTWINS_EDGE_BIND:
		msg.Router.Target = "edge"
		//send to event hub.
		dtc.SendToModule("EventHub", msg)
		//cache the message.
		dtm.dtcontext.CacheMessage(msg)	
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
		//TODO:

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
		twinMsg, err := common.UnMarshalTwinMessage(msg)
	if err != nil {
		return
	}	
	dgTwin := &twinMsg.Twins[0]
	case common.DGTWINS_OPS_DELETE:
	case common.DGTWINS_OPS_GET:
	case common.DGTWINS_OPS_WATCH
	}

	//if send success, the cache the message.
	
}
