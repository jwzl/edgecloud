package config

import (
	"k8s.io/klog"
	"github.com/jwzl/beehive/pkg/common/config"
)

type ServerConfig struct {
	Port	string
}

func GetServerConfig() *ServerConfig {
	setting := &ServerConfig{}

	Port, err := config.CONFIG.GetValue("web.server.port").ToString()
	if err != nil {
		klog.Errorf("Failed to get broker url for mqtt client: %v", err)
		Port = "http://127.0.0.1"
	}
	setting.Port = Port

	return setting
}
