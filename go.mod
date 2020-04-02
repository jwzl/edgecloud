module github.com/jwzl/edgecloud

go 1.13

require (
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.6.2
	github.com/jwzl/beehive v0.0.0-20191028085830-1606e1f5c86a
	github.com/jwzl/edgeOn v1.2.0
	github.com/jwzl/mqtt v1.3.0
	github.com/jwzl/wssocket v1.0.0
	github.com/spf13/cobra v0.0.6
	k8s.io/component-base v0.17.4
	k8s.io/klog v1.0.0
)

replace github.com/jwzl/mqtt v1.3.0 => github.com/jwzl/mqtt v0.0.0-20200310015455-e512045b0629

replace github.com/jwzl/edgeOn v1.2.0 => github.com/jwzl/edgeOn v0.0.0-20200402083535-9ef0b8b7ff25
