package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jwzl/edgecloud/pkg/web/apis"
)

func InitEdgeRouter(Router *gin.RouterGroup) {
	EdgeRouter := Router.Group("edge")
	{
		EdgeRouter.POST("/bind", apis.BindEdgeApi)
		EdgeRouter.GET("/list", apis.ListEdgeApi)
	}
}
