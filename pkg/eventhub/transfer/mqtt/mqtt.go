package mqtt

import (
	"fmt"
	"sync"
	"time"
	"strings"
	"k8s.io/klog"
	"github.com/jwzl/mqtt/client"
	"github.com/jwzl/wssocket/fifo"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/edgeOn/msghub/config"
)


const (
	retryCount       = 5
	cloudAccessSleep = 5 * time.Second

	//mqtt topic should has the format:
	// mqtt/dgtwin/cloud[edge]/{edgeID}/comm for communication.
	// mqtt/dgtwin/cloud[edge]/{edgeID}/control  for some control message.
	MQTT_SUBTOPIC_PREFIX	= "mqtt/dgtwin/cloud"
	MQTT_PUBTOPIC_PREFIX	= "mqtt/dgtwin/edge"
)

type Client	struct {
	Settings	*config.MqttConfig

	subTopics	[]string
	// for mqtt send thread.
	mutex 		sync.RWMutex
	mqttClient		*client.MQTTClient
}

func NewMqttClient(setting *config.MqttConfig) *Client {
	mc := &Client{
		Settings: setting,
		subTopics: make([]string, 0),		
	}

	c := &client.MQTTClient{
		Host:			setting.URL,
		User:			setting.User,
		Passwd:			setting.Passwd,	
		ClientID:		setting.ClientID,
		keepAliveInterval:	setting.keepAliveInterval,
		PingTimeout:		setting.PingTimeout,
		MessageChannelDepth: setting.MessageCacheDepth,			
		CleanSession:	true,
		FileStorePath: "memory",
		OnConnect:	mc.ClientOnConnect,
		OnLost:		mc.ClientOnLost,			
		WillTopic:		"",			//no will topic.	
		TLSConfig:		setting.tlsConfig,  

	}
	mc.mqttClient = c

	return mc
}

func (mc *Client) Start(){

	//Start the client
	mc.mqttClient.Start()

	// retry to connect.
	err := mc.tryToConnect()
	if err != nil {
		klog.Errorf("Client connecte err (%v), Please check your net link", err)	
		return
	}
}

func (mc *Client) tryToConnect() error {
	for i := 0; i < retryCount; i++ {
		err := mc.mqttClient.Connect()
		if err != nil {
			klog.Errorf("Client connecte err (%v), retry...", err)
		}else {
			klog.Infof("Client connecte Successful")
			return nil
		}
		time.Sleep(cloudAccessSleep)
	}

	return err
}

func (mc *Client) messageDispatch(topic string, msg *model.Message){

}

//ClientOnConnect
func (mc *Client) ClientOnConnect(client mqtt.Client) {
	klog.Infof("Connect to mqtt broker %s Successful.", mc.setting.URL)
}

//ClientOnLost
func (mc *Client) ClientOnLost(client mqtt.Client, err error) {
	klog.Infof("MQTT connection is lost, we restart this.")
	//Restart the mqtt client
	mc.Close()
	mc.Start()
}

/*
* publish the message.
*/
func (mc *Client) Publish(topic string, msg *model.Message) error {
	return mc.mqttClient.Publish(topic, mc.Settings.QOS, mc.Settings.Retain, msg)
}

/*
* Subscribe the message.
*/
func (mc *Client) Subscribe(topic string) error {
	return mc.mqttClient.Subscribe(topic, mc.Settings.QOS, mc.messageDispatch)
}

/*
* UnSubscribe the message.
*/
func (mc *Client) Unsubscribe(topics string) error {
	return mc.mqttClient.Unsubscribe(topics)
}

/*
* Close the mqtt client.
*/
func (mc *Client) Close(){
	mc.mqttClient.Close()
}
