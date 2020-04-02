package  web

import (
	"time"
	"net/http"
	"k8s.io/klog"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
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
	ApiGroup := Router.Group("rest/v1")
	router.InitApiRouter(ApiGroup)
	router.InitDeviceRouter(ApiGroup)
	router.InitEdgeRouter(ApiGroup)

	return Router
}

func Cors() gin.HandlerFunc {
   return func(c *gin.Context) {
      	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		klog.Infof("XXXXXXXXX",)
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(200)
            return
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
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "OPTIONS", "POST", "GET"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
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
