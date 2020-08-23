package devicetwin

import (
	"time"
	"sync"
	"errors"
	"strings"
	"k8s.io/klog"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/edgecloud/pkg/types"
	"github.com/jwzl/beehive/pkg/core/context"
	"github.com/jwzl/edgecloud/pkg/devicetwin/eventlistener"
)


type DTContext struct {
	Context			*context.Context
	/*
	* This is all edge in this cluster.
	*/
	EdgeMap			*sync.Map
	EdgeMutex		*sync.Map
	MessageCache	*sync.Map
	EdgeHealth		*sync.Map
}

func NewDTContext(c *context.Context) *DTContext {
	if c == nil {
		return nil
	}

	var edges sync.Map
	var edgesMutex sync.Map
	var cache sync.Map
	var healthEdge sync.Map

	return &DTContext{
		Context:	c,
		EdgeMap:	&edges,
		EdgeMutex:  &edgesMutex,
		MessageCache: &cache,
		EdgeHealth: &healthEdge,
	}
}


//GetEdgeMutex get the edge mutex
func (dtc *DTContext) GetEdgeMutex (edgeID string) (*sync.Mutex, bool) {
	v, exist := dtc.EdgeMutex.Load(edgeID)
	if !exist {
		return nil, false
	}

	mutex, isMutex := v.(*sync.Mutex)
	if !isMutex {
		return nil, false
	}

	return mutex, true
}

//Lock twin by ID
func (dtc *DTContext) Lock (edgeID string) bool {
	mutex, ok := dtc.GetEdgeMutex(edgeID)
	if ok {
		mutex.Lock()
		return true
	}

	return false
}

//unlock twin by ID 
func (dtc *DTContext) Unlock (edgeID string) bool {
	mutex, ok := dtc.GetEdgeMutex(edgeID)
	if ok {
		mutex.Unlock()
		return true
	}

	return false
}

/*
* GetEdgeInfo get the edge information.
*/
func (dtc *DTContext) GetEdgeInfo(edgeID string) (*EdgeDescription, error){
	v, exist := dtc.EdgeMap.Load(edgeID)
	if !exist {
		return nil, errors.New("edge is not exists.")
	}

	edged, isEdged := v.(*EdgeDescription)
	if !isEdged {
		return nil,  errors.New("edge is not exists.")
	}

	return edged, nil
} 

/*
* AddEdgeInfo: Add edge information.
*/
func (dtc *DTContext) AddEdgeInfo(edged *EdgeDescription) error {
	edgeID := edged.ID
	if _, err := dtc.GetEdgeInfo(edgeID); err == nil {
		return errors.New("edge is exists.")
	}

	var edgeMutex	sync.Mutex
	dtc.EdgeMutex.Store(edgeID, &edgeMutex)
	dtc.EdgeMap.Store(edgeID, edged)

	//notify the edge register event.
	eventlistener.MatchEventAndDispatch(edgeID,
		"", eventlistener.EVENT_EDGE_CREATED)

	return nil
}

/*
* List edge info.
*/
func (dtc *DTContext) ListEdgeInfo() []types.EdgeInfo {
	edgeds := make([]types.EdgeInfo, 0)
	
	dtc.EdgeMap.Range(func(_, value interface{}) bool{
		edged := value.(*EdgeDescription)
		if edged == nil {
			return true
		}
		
		edgeds = append(edgeds, types.EdgeInfo{
			ID: edged.ID,
			Name: edged.Name,
			Description: edged.Description,
			State: 	edged.State,
			DeviceIDs: edged.GetDeviceIDs(),
		})
		
		return true
	})
	
	return edgeds
}
/*
* Edge is online?
*/
func (dtc *DTContext) EdgeIsOnline(edgeID string) bool {
	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return false
	}
	
	if edged.GetEdgeState() != EdgeStateOnline {
		return false
	}

	return true
}
/*
* Set edge State
*/
func (dtc *DTContext) SetEdgeState(edgeID, state string) error {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)

	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return err
	}
	
	if edged.GetEdgeState() != state {
		edged.SetEdgeState(state)
	}

	return nil
}

/*
* RegisterTwins: register the twin. 
*/
func (dtc *DTContext) RegisterTwins(edgeID, twinID string) error {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)

	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return err
	}

	return edged.RegisterTwins(twinID)
}

/*
* UnRegisterTwins: Unregister the twin. 
*/
func (dtc *DTContext) UnRegisterTwins(edgeID, twinID string) error {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)

	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return err
	}

	return edged.UnRegisterTwins(twinID)
}

/*
* UpdateTwin update twin in this edge
*/
func (dtc *DTContext) UpdateTwin(edgeID string, twin *common.DigitalTwin) error {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)

	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return err
	}

	return edged.UpdateTwin(twin)
}

/*
* DeleteTwin delete the twin.
*/
func (dtc *DTContext) DeleteTwin(edgeID, twinID string) error {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)

	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return err
	}

	return edged.DeleteTwin(twinID)
}

/*
* GetRawTwin: get the raw twin. 
*/
func (dtc *DTContext) GetRawTwin(edgeID, twinID string) ([]byte, error) {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)

	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return nil, err
	}

	return edged.GetRawTwin(twinID)
}

func (dtc *DTContext) ListTwins (edgeID string) ([]common.DigitalTwin, error) {
	dtc.Lock(edgeID)
	defer dtc.Unlock(edgeID)
	
	edged, err := dtc.GetEdgeInfo(edgeID)
	if err != nil {
		return nil, err
	}
	
	return edged.ListTwins(), nil
}

//SendResponseMessage Send Response conten.
func (dtc *DTContext) SendResponseMessage(requestMsg *model.Message, content []byte){
	resource := requestMsg.GetResource()
	target := requestMsg.GetSource()	

	modelMsg := common.BuildModelMessage(common.CloudName, target, 
					common.DGTWINS_OPS_RESPONSE, resource, content)	
	modelMsg.SetTag(requestMsg.GetID())	
	klog.Infof("Send response message (%v)", modelMsg)

	dtc.Context.Send(types.EDGECLOUD_EVENTHUB_MODULE, *modelMsg)
}

func (dtc *DTContext) SendTwinMessage(edgeID, operation string, content []byte){
	resource := edgeID+"/"+common.DGTWINS_RESOURCE_TWINS

	modelMsg := common.BuildModelMessage(common.CloudName, common.TwinModuleName, 
					operation, resource, content)	

	dtc.Context.Send(types.EDGECLOUD_EVENTHUB_MODULE, *modelMsg)
	//cache the message.
	dtc.CacheMessage(modelMsg)
}

func (dtc *DTContext) SendPropertyMessage(edgeID, operation string, content interface{}){
	resource := edgeID+"/"+common.DGTWINS_RESOURCE_PROPERTY

	modelMsg := common.BuildModelMessage(common.CloudName, common.TwinModuleName, 
					operation, resource, content)	

	dtc.Context.Send(types.EDGECLOUD_EVENTHUB_MODULE, *modelMsg)
	//cache the message.
	dtc.CacheMessage(modelMsg)
}

func (dtc *DTContext) CacheMessage(msg *model.Message){
	msgID := msg.GetID()
 
	_, exist := dtc.MessageCache.Load(msgID)
	if !exist {
		dtc.MessageCache.Store(msgID, msg)
	}
}

func (dtc *DTContext) CacheHasThisMessage(msg *model.Message) bool {
	msgID := msg.GetID()
 
	_, exist := dtc.MessageCache.Load(msgID)
	return exist
}

func (dtc *DTContext) DeleteMsgCache(msg *model.Message){
	pMsgID := msg.GetTag()
 
	if pMsgID != "" {
		_, exist := dtc.MessageCache.Load(pMsgID)
		if exist {
			dtc.MessageCache.Delete(pMsgID) 
		}
	}
}

//dealHeartBeat deal heartbeat message from edge.
func (dtc *DTContext) DealHeartBeat(msg *model.Message) {
	resource  := msg.GetResource()
	splitString := strings.Split(resource, "/")	
	edgeID := splitString[0]

	dtc.EdgeHealth.Store(edgeID, time.Now().Unix())
}
