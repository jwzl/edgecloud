package main


import (
	"k8s.io/klog"
	"github.com/jwzl/edgecloud/pkg/eventhub/settings"
)

func main() {
	_, err := settings.GetMqttSetting()
	if err != nil {
		klog.Infof("get config with err %v", err)
	}

}
