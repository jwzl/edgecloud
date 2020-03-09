package transfer

import (
	"github.com/jwzl/edgecloud/pkg/eventhub/settings"
	"github.com/jwzl/edgecloud/pkg/eventhub/transfer/mqtt"
)
 
/*
* Transfer interfaces
*/
type Transfer interface {
	//Start the connection 
	Start()
	//publish the message
	Publish(string, *model.Message) error
	// subscribe the message;
	Subscribe(topic string, fn func(string, *model.Message)) error
	//unsubscribe the message;
	Unsubscribe(string) error
	// close the connection.
	Close()
}


func NewTransfer(setting *settings.MqttSettings) Transfer {
	return mqtt.NewMqttClient(setting)
}
