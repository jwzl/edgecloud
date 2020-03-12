package devicetwin

import (
	"time"
	"sync"
	"errors"
	"strings"
	"k8s.io/klog"
	"github.com/jwzl/edgeOn/common"
	"github.com/jwzl/wssocket/model"
	"github.com/jwzl/beehive/pkg/core/context"
)


type DTContext struct {
	Context			*context.Context
	/*
	* This is all edge in this cluster.
	*/
	EdgeMap		*sync.Map
	EdgeMutex	*sync.Map
}

func NewDTContext(c *context.Context) *DTContext {
	if c == nil {
		return nil
	}

	var edges sync.Map
	var edgesMutex sync.Map

	return &DTContext{
		Context:	c,
		EdgeMap:	&edges,
		EdgeMutex:  &edgesMutex,
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
	edgeID : = edged.ID
	if err := dtc.GetEdgeInfo(edgeID); err == nil {
		return errors.New("edge is exists.")
	}

	var edgeMutex	sync.Mutex
	dtc.EdgeMutex.Store(edgeID, &edgeMutex)
	dtc.EdgeMap.Store(edgeID, edged)

	return nil
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
	
	edged.SetEdgeState(state)

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
		return err
	}

	return edged.GetRawTwin(twinID)
}
