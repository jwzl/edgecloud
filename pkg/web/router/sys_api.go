package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jwzl/edgecloud/pkg/web/apis"
)

func InitApiRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("api")
	{
		ApiRouter.GET("/users", apis.GetUsersHandler)
		ApiRouter.POST("/user", apis.AddUserHandler)
		ApiRouter.PUT("/user/:id", apis.UpdateUserHandler)
		ApiRouter.DELETE("/user/:id", apis.DeleteUserHandler)
	}
}
