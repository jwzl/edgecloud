package devicetwin


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

		if strings.Contains(target, common.CloudName) {
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
		
	case common.DGTWINS_OPS_SYNC:
		twinMsg, err := common.UnMarshalTwinMessage(msg)
		if err != nil {
			return
		}	

		dgTwin := &twinMsg.Twins[0]

		err := dtm.dtcontext.UpdateTwin(edgeID, dgTwin)
		if err != nil {
			return
		}	
		klog.Infof("Update successful.") 
	}	
}

/*
* doDownStreamMessage
*/
func (dtm *DeviceTwinModule) doDownStreamMessage(msg *model.Message) {
	operation := msg.GetOperation()

	switch operation {
	case common.DGTWINS_OPS_RESPONSE:
		
	case common.DGTWINS_OPS_SYNC:
				
	}	
}
