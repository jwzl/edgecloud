package mqtt

import (
	"sync"
	"time"
	"k8s.io/klog"
	"github.com/jwzl/mqtt/client"
	"github.com/jwzl/wssocket/model"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jwzl/edgecloud/pkg/eventhub/settings"
)


const (
	retryCount       = 5
	cloudAccessSleep = 5 * time.Second

	//mqtt topic should has the format:
	// mqtt/dgtwin/cloud[edge]/{edgeID}/comm for communication.
	// mqtt/dgtwin/cloud[edge]/{edgeID}/control  for some control message.
	//MQTT_SUBTOPIC_PREFIX	= "mqtt/dgtwin/cloud"
	//MQTT_PUBTOPIC_PREFIX	= "mqtt/dgtwin/edge"
)

type Client	struct {
	Settings	*settings.MqttSettings
	// for mqtt send thread.
	mutex 		sync.RWMutex
	mqttClient		*client.MQTTClient
}

func NewMqttClient(setting *settings.MqttSettings) *Client {
	mc := &Client{
		Settings: setting,	
	}

	c := &client.MQTTClient{
		Host:			setting.URL,
		User:			setting.User,
		Passwd:			setting.Passwd,	
		ClientID:		setting.ClientID,
		Order:			true,
		KeepAliveInterval:	setting.KeepAliveInterval,
		PingTimeout:		setting.PingTimeout,
		MessageChannelDepth: setting.MessageCacheDepth,			
		CleanSession:	true,
		FileStorePath: "memory",
		OnConnect:	mc.ClientOnConnect,
		OnLost:		mc.ClientOnLost,			
		WillTopic:		"",			//no will topic.	
		TLSConfig:		setting.TLSConfig,  

	}
	mc.mqttClient = c

	return mc
}

func (mc *Client) Start(){

	//Start the client
	mc.mqttClient.Start()

	// retry to connect.
retry_connect:
	err := mc.tryToConnect()
	if err != nil {
		klog.Errorf("Client connecte err (%v), Please check your net link", err)	
		goto retry_connect
	}
}

func (mc *Client) tryToConnect() error {
	var err error

	for i := 0; i < retryCount; i++ {
		err = mc.mqttClient.Connect()
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

//ClientOnConnect
func (mc *Client) ClientOnConnect(client mqtt.Client) {
	klog.Infof("Connect to mqtt broker %s Successful.", mc.Settings.URL)
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
	mc.mutex.Lock()
	defer mc.mutex.Unlock()
	
	return mc.mqttClient.Publish(topic, mc.Settings.QOS, mc.Settings.Retain, msg)
}

/*
* Subscribe the message.
*/
func (mc *Client) Subscribe(topic string, fn func(topic string, msg *model.Message)) error {
	return mc.mqttClient.Subscribe(topic, mc.Settings.QOS, fn)
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
