/*
* event hub:
*/

package eventhub

import(
	"github.com/jwzl/beehive/pkg/core"
	"github.com/jwzl/beehive/pkg/core/context"
)

type EventHub struct {
	trans	transfer.Transfer
}

// Register this module.
func Register(){	
	setting := 
	eh := &EventHub{
		trans : transfer.NewTransfer(setting),
	}
	core.Register(eh)
}

//Name
func (eh *EventHub) Name() string {
	return "EventHub"
}

//Group
func (eh *EventHub) Group() string {
	return "EventHub"
}

//Start this module.
func (eh *EventHub) Start(c *context.Context) {
	klog.Infof("Start the module!")


	for {
		v, err := c.Receive(eh.Name())
		if err != nil {
			klog.Errorf("failed to receive message: %v", err)
			break
		}

		msg, isThisType := v.(*model.Message)
		if !isThisType || msg == nil { 		
			continue
		}

		target := msg.GetTarget()
		switch target {
		case "edge":
			operation := msg.GetOperation()
			switch operation {
			case common.DGTWINS_OPS_CREATE:

			case "KeepAlive":
					
			}	
		case common.TwinModuleName:
			 
		}
	}
}

//Cleanup
func (eh *EventHub) Cleanup() {
	eh.context.Cleanup(eh.Name())
}
