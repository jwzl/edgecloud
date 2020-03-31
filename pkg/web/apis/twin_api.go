package apis

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

//POST /edge/bind?edgeid=xxx
func BindEdgeApi(c *gin.Context) {
	edgeId := c.Query("edgeid")
	if edgeId == "" {
		c.JSON(http.StatusBadRequest,  "edgeid is empty")
		return 
	}
	if err := BindEdge(edgeId); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, nil);
	}
}

//POST /create?edgeID=001&twinID=001
func CreateTwinApi(c *gin.Context)  {
	edgeId := c.Query("edgeid")
	twinId := c.Query("twinid")
	if edgeId == "" || twinId == "" {
		c.JSON(http.StatusBadRequest, "edgeid or twinid is empty")
		return 
	}

	if err := CreateTwin(edgeId, twinId); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, nil);
	}
}

//DELETE /delete?edgeID=001&twinID=001
func DeleteTwinApi(c *gin.Context) {
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

//Get /get?edgeID=001&twinID=001
func GetTwinApi(c *gin.Context) {
	edgeId := c.Query("edgeid")
	twinId := c.Query("twinid")
	if edgeId == "" || twinId == "" {
		c.JSON(http.StatusBadRequest, "edgeid or twinid is empty")
		return 
	}

	if dgTwin, err := GetTwin(edgeId, twinId); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return 
	} else {
		c.JSON(http.StatusOK, dgTwin);
	}
}
