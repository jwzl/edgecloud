package  web

import (
	"time"
	"net/http"
	"k8s.io/klog"
	"github.com/gin-gonic/gin"
	"github.com/jwzl/beehive/pkg/core"
	"github.com/jwzl/edgecloud/pkg/types"
	"github.com/jwzl/edgecloud/pkg/web/apis"
	"github.com/jwzl/edgecloud/pkg/web/router"
	"github.com/jwzl/edgecloud/pkg/web/config"
	"github.com/jwzl/beehive/pkg/core/context"	
)

type WebModule struct {
	context *context.Context	
}

/*
* Init router
*/
func InitRouter() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("")
	router.InitApiRouter(ApiGroup)
	router.InitDeviceRouter(ApiGroup)
	router.InitEdgeRouter(ApiGroup)

	return Router
}

func Cors() gin.HandlerFunc {
   return func(c *gin.Context) {
      method := c.Request.Method

      c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
      c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
      c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
      c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
      c.Header("Access-Control-Allow-Credentials", "true")

      //allow OPTIONS method
      if method == "OPTIONS" {
         c.AbortWithStatus(http.StatusNoContent)
      }
      c.Next()
   }
}

// Register this module.
func Register(){
	wm := &WebModule{}
	core.Register(wm)
}

//Name
func (wm *WebModule) Name() string {
	return types.EDGECLOUD_APISERVER_MODULE
}

//Group
func (wm *WebModule) Group() string {
	return types.EDGECLOUD_APISERVER_MODULE
}

//Start this module.
func (wm *WebModule) Start(c *context.Context) {
	wm.context = c
	apis.NewDeviceTwinModule(c)
	//init router
	router := InitRouter()  
	router.Use(Cors())
	port := config.GetServerConfig().Port

	s := &http.Server{
		Addr:              port,
		Handler:           router,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}

	klog.Infof("[Info] start http server listening %s", port)
restart:
	err := s.ListenAndServe()
	if err != nil {
		klog.Infof(" http ListenAndServe with err:%s", err)
		time.Sleep(3 * time.Second)
		goto restart
	}
}
//Cleanup
func (wm *WebModule) Cleanup() {
	wm.context.Cleanup(wm.Name())
}
