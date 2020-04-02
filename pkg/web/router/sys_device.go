package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jwzl/edgecloud/pkg/web/apis"
)

func InitDeviceRouter(Router *gin.RouterGroup) {
	DeviceRouter := Router.Group("dev")
	{
		DeviceRouter.GET("/get", apis.GetTwinApi)
		DeviceRouter.POST("/twin", apis.CreateTwinApi)
		DeviceRouter.DELETE("/delete", apis.DeleteTwinApi)
	}
}
