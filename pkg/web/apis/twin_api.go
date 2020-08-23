package apis

import (
	"net/http"
	_"k8s.io/klog"
	"github.com/gin-gonic/gin"
)

func AddCors(c *gin.Context){
	/*
	* Gin post request can't call gin.Use() api correctlly.
	* We add this tag into header
	*/
	header := c.Writer.Header()
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
}

//POST /edge/bind?edgeid=xxx
func BindEdgeApi(c *gin.Context) {
	AddCors(c)
	edgeId := c.Query("edgeid")
	if edgeId == "" {
		c.JSON(http.StatusBadRequest,  "edgeid is empty")
		return 
	}

	c.JSON(BindEdge(edgeId))
}

//POST /create?edgeID=001&twinID=001
func CreateTwinApi(c *gin.Context)  {
	AddCors(c)
	edgeId := c.Query("edgeid")
	twinId := c.Query("twinid")
	if edgeId == "" || twinId == "" {
		c.JSON(http.StatusBadRequest, "edgeid or twinid is empty")
		return 
	}

	c.JSON(CreateTwin(edgeId, twinId))
}

//DELETE /delete?edgeID=001&twinID=001
func DeleteTwinApi(c *gin.Context) {
	AddCors(c)
	edgeId := c.Query("edgeid")
	twinId := c.Query("twinid")
	if edgeId == "" || twinId == "" {
		c.JSON(http.StatusBadRequest, "edgeid or twinid is empty")
		return 
	}

	if err := DeleteTwin(edgeId, twinId); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, nil);
	}
}

//Get /edge/twin?edgeID=001&twinID=001
func GetTwinApi(c *gin.Context) {
	AddCors(c)
	edgeId := c.Query("edgeid")
	twinId := c.Query("twinid")
	if edgeId == "" || twinId == "" {
		c.JSON(http.StatusBadRequest, "edgeid or twinid is empty")
		return 
	}

	if twins, err := GetTwin(edgeId, twinId); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, twins);
	}
}

//Get /dev/list?edgeID=001
func ListTwinApi(c *gin.Context) {
	AddCors(c)
	edgeId := c.Query("edgeid")
	if edgeId == "" {
		c.JSON(http.StatusBadRequest, "edgeid is empty")
		return 
	}
	if twins, err := ListTwins(edgeId); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, twins);
	}
	
}

//Get /edge/list
func ListEdgeApi(c *gin.Context) {
	AddCors(c)
	
	if edges, err := ListEdge(); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, edges);
	}	
}		
