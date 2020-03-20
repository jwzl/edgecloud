package main


import (
	"k8s.io/klog"
	"k8s.io/component-base/logs"
	"github.com/jwzl/edgecloud/cmd"
)

func main() {
	command := cmd.NewAppCommand()
	logs.InitLogs()
	defer logs.FlushLogs()	

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
