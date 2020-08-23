package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jwzl/edgecloud/pkg/web/apis"
)

func InitDeviceRouter(Router *gin.RouterGroup) {
	DeviceRouter := Router.Group("dev")
	{
		DeviceRouter.GET("/twin", apis.GetTwinApi)
		DeviceRouter.GET("/list", apis.ListTwinApi)	
		DeviceRouter.PUT("/twin", apis.CreateTwinApi)
		DeviceRouter.DELETE("/delete", apis.DeleteTwinApi)
	}
}
