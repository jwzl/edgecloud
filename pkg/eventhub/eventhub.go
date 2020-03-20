/*
* event hub:
*/
package eventhub

import(
	"fmt"
	"strings"
	"k8s.io/klog"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/jwzl/beehive/pkg/core/context"
	"github.com/jwzl/edgecloud/pkg/eventhub/settings"
	"github.com/jwzl/edgecloud/pkg/eventhub/transfer"
)

const (
	//mqtt topic should has the format:
	// mqtt/dgtwin/cloud[edge]/{edgeID}/comm for communication.
	// mqtt/dgtwin/cloud[edge]/{edgeID}/control  for some control message.
	MQTT_SUBTOPIC_PREFIX	= "mqtt/dgtwin/cloud"
	MQTT_PUBTOPIC_PREFIX	= "mqtt/dgtwin/edge"
)
type EventHub struct {
	trans	transfer.Transfer
	context *context.Context
}

// Register this module.
func Register(){	
	setting, err := settings.GetMqttSetting()
	if err != nil {
		klog.Errorf("Get mqtt setting failed")
		return
	}

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
	var topic string
	klog.Infof("Start the module!")
	
	eh.context = c
	//Start the client.
	eh.trans.Start()
	topic = fmt.Sprintf("%s/#", MQTT_SUBTOPIC_PREFIX)	
	eh.trans.Subscribe(topic, eh.messageDispatch)

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

		resource := msg.GetResource()
		splitString := strings.Split(resource, "/")
		edgeID := splitString[0]
		resource = splitString[1]
		msg.Router.Resource = resource
		target := msg.GetTarget()

		switch target {
		case "edge":
			operation := msg.GetOperation()
			if operation == "Bind" {
				/*
				* this will bind the edge by edgeID at remote. After that, 
				* the edge will send heartbeat to us.
				*/		
				topic = fmt.Sprintf("%s/%s/bind", MQTT_PUBTOPIC_PREFIX, edgeID)
			}
		case common.TwinModuleName:
			topic = fmt.Sprintf("%s/%s/comm", MQTT_PUBTOPIC_PREFIX, edgeID)	
		}

		err = eh.trans.Publish(topic, msg)
		if err != nil {
			klog.Errorf("failed to publish message: %v", err)
			return
		}
	}
}

//Cleanup
func (eh *EventHub) Cleanup() {
	eh.context.Cleanup(eh.Name())
}

func (eh *EventHub) messageDispatch (topic string, msg *model.Message){
	splitString := strings.Split(topic, "/")

	if len(splitString) < 5 {
		//ignore 
		return
	}

	edgeID := splitString[3]
	resource := edgeID+"/"+msg.GetResource()
	msg.Router.Resource = resource
	
	if splitString[4]  == "hearbeat" {
		eh.context.Send("deviceTwin", msg)	
	}else if splitString[4]  == "comm" {
		//send to device twin module.
		eh.context.Send("deviceTwin", msg)
	}
}
