package devicetwin


type DeviceTwinModule struct {
	context 	*context.Context
	dtcontext 	*DTContext
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

	for {
		v, err := c.Receive(dtm.Name())
		if err != nil {
			klog.Errorf("failed to receive message: %v", err)
			break
		}

		msg, isThisType := v.(*model.Message)
		if !isThisType || msg == nil { 		
			continue
		}

		operation := msg.GetOperation()
		switch operation {
		case common.DGTWINS_OPS_RESPONSE:
		
		case common.DGTWINS_OPS_SYNC:	
		}	
	}
}

//Cleanup
func (dtm *DeviceTwinModule) Cleanup() {
	dtm.context.Cleanup(dtm.Name())
}
