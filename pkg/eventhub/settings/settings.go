package settings

import (
	"time"
	"crypto/tls"
	"k8s.io/klog"
	"github.com/jwzl/beehive/pkg/common/config"
)

type MqttSettings struct {
	URL				string
	ClientID		string
	User			string
	Passwd			string
	// tls config
	TLSConfig 		*tls.Config
	KeepAliveInterval	time.Duration
	PingTimeout			time.Duration  
	QOS				 	byte
	Retain			   	bool	
	MessageCacheDepth  	uint
}

func GetMqttSetting()(*MqttSettings, error){
	setting := &MqttSettings{}

	url, err := config.CONFIG.GetValue("eventhub.mqtt.broker").ToString()
	if err != nil {
		klog.Errorf("Failed to get broker url for mqtt client: %v", err)
		return nil, err
	}
	setting.URL = url

	id, err := config.CONFIG.GetValue("eventhub.mqtt.clientid").ToString()
	if err != nil {
		klog.Warningf("Failed to get client id: %v", err)
		return nil, err
	}
	setting.ClientID = id

	user, err := config.CONFIG.GetValue("eventhub.mqtt.user").ToString()
	if err != nil {
		klog.Infof("eventhub.mqtt.user is empty")
		user = ""
	}
	setting.User = user

	passwd, err := config.CONFIG.GetValue("eventhub.mqtt.passwd").ToString()
	if err != nil {
		klog.Infof("eventhub.mqtt.passwd is empty")
		passwd = ""
	}
	setting.Passwd = passwd	

	certfile, err := config.CONFIG.GetValue("eventhub.mqtt.certfile").ToString()
	if err != nil {
		klog.Infof("eventhub.mqtt.certfile is empty")
		certfile = ""
	}

	keyfile, err := config.CONFIG.GetValue("eventhub.mqtt.keyfile").ToString()
	if err != nil {
		klog.Infof("eventhub.mqtt.keyfile is empty")
		keyfile = ""
	}
	if certfile != "" && keyfile != "" {
		tlsConfig, err := CreateTLSConfig(certfile, keyfile)
		if err != nil {
			klog.Infof("TLSConfig Disabled")
			setting.TLSConfig = nil
		}else{
			setting.TLSConfig = tlsConfig
		}	
	}
	
	keepAliveInterval, err := config.CONFIG.GetValue("eventhub.mqtt.keep-alive-interval").ToInt()
	if err != nil {
		klog.Infof("eventhub.mqtt.keep-alive-interval is empty")
		keepAliveInterval = 120
	}
	setting.KeepAliveInterval = time.Duration(keepAliveInterval) * time.Second

	pingTimeout, err := config.CONFIG.GetValue("eventhub.mqtt.ping-timeout").ToInt()
	if err != nil {
		klog.Infof("eventhub.mqtt.ping-timeout is empty")
		pingTimeout = 120
	}
	setting.PingTimeout = time.Duration(pingTimeout) * time.Second

	qos, err := config.CONFIG.GetValue("eventhub.mqtt.qos").ToInt()
	if err != nil {
		klog.Infof("eventhub.mqtt.qos is empty")
		qos = 2
	}
	setting.QOS = uint8(qos)

	retain, err := config.CONFIG.GetValue("eventhub.mqtt.retain").ToBool()
	if err != nil {
		klog.Infof("eventhub.mqtt.retain is empty")
		retain = false
	}
	setting.Retain = retain

	sessionQueueSize, err := config.CONFIG.GetValue("eventhub.mqtt.session-queue-size").ToInt()
	if err != nil {
		klog.Infof("eventhub.mqtt.session-queue-size is empty")
		sessionQueueSize = 100
	}
	setting.MessageCacheDepth = uint(sessionQueueSize)

	return setting, nil
}

func CreateTLSConfig(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	return tlsConfig, nil
}
