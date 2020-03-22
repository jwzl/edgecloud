/*
* The edgemgr has some function as below:
* 1. Detect  the edge cluster on remote.
* 2. Keepalive between cloud and edge sides.
* 3. Manage the edge on remote.
*/
package devicetwin


import (
	"sync"
	"errors"

	"encoding/json"
	"github.com/jwzl/edgeOn/common"
)

const (
	EdgeStateInitial	= "initial"
	EdgeStateOnline		= "online"
	EdgeStateOffline	= "offline"
)
type EdgeDescription struct {
	ID				string
	Name 			string
	Description		string
	State			string
	/*
	* all device ID in this edge.
	*/
	deviceIDs  		[]string
	/*
	* all device twin in this edge;
	*/
	Twins			*sync.Map
	TwinMutex		*sync.Map
}


func NewEdgeDescription(edgeID string) *EdgeDescription {
	var twins		sync.Map
	var twinMutex	sync.Map

	return &EdgeDescription{
		ID: edgeID,
		State: EdgeStateInitial,
		deviceIDs: make([]string, 0),
		Twins: &twins,
		TwinMutex: &twinMutex, 
	}
}


func (ed *EdgeDescription) SetEdgeName(name string) {
	ed.Name = name
}

func (ed *EdgeDescription) SetEdgeDescription(desp string) {
	ed.Description = desp
}

func (ed *EdgeDescription) SetEdgeState(state string) {
	ed.State = state
}

func (ed *EdgeDescription) FindTwins(twinID string) bool {
	for _, ids := range ed.deviceIDs {
		if ids == twinID {
			return true
		}
	}

	return false
}

func (ed *EdgeDescription) RegisterTwins(twinID string) error {
	if ed.FindTwins(twinID) != true {
		ed.deviceIDs = append(ed.deviceIDs, twinID)
		return nil
	}

	return errors.New("twin is exists.")
}

func (ed *EdgeDescription) UnRegisterTwins(twinID string) error {
	if ed.FindTwins(twinID) != true {
		return errors.New("twin is not exists.")
	}

	for key, ids := range ed.deviceIDs {
		if ids == twinID {
			ed.deviceIDs = append(ed.deviceIDs[:key], ed.deviceIDs[key+1:]...)
		}
	}
	
	return ed.DeleteTwin(twinID)
}

func (ed *EdgeDescription) getTwin (twinID string) (*common.DigitalTwin, bool) {
	v, exist := ed.Twins.Load(twinID)
	if !exist {
		return nil, false
	}

	twin, isTwin := v.(*common.DigitalTwin)
	if !isTwin {
		return nil, false
	}

	return twin, true
}

//getTwinMutex get the twin mutex
func (ed *EdgeDescription) getTwinMutex (twinID string) (*sync.Mutex, bool) {
	v, exist := ed.TwinMutex.Load(twinID)
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
func (ed *EdgeDescription) Lock (twinID string) bool {
	mutex, ok := ed.getTwinMutex(twinID)
	if ok {
		mutex.Lock()
		return true
	}

	return false
}

//unlock twin by ID 
func (ed *EdgeDescription) Unlock (twinID string) bool {
	mutex, ok := ed.getTwinMutex(twinID)
	if ok {
		mutex.Unlock()
		return true
	}

	return false
}


/*
* Update twin.
*/
func (ed *EdgeDescription) UpdateTwin(twin *common.DigitalTwin) error {
	twinID := twin.ID
	if ed.FindTwins(twinID) != true {
		return errors.New("No such this twin")
	}

	_ , exist := ed.TwinMutex.Load(twinID)
	if !exist {
		//Create twin
		var deviceMutex	sync.Mutex
		ed.TwinMutex.Store(twinID, &deviceMutex)

		ed.Twins.Store(twinID, twin)
	} else {
		// update twin
		ed.Lock(twinID)
		defer ed.Unlock(twinID)	
	
		oldTwin, ok := ed.getTwin(twinID) 
		if !ok {
			return errors.New("No such twin")
		}
		//update the twin
		if len(twin.Name) > 0 {
			oldTwin.Name = twin.Name
		}
		if len(twin.Description) > 0 {
			oldTwin.Description = twin.Description
		}
		if len(twin.State) > 0 {
			oldTwin.LastState = twin.State 
			oldTwin.State = twin.State		
		}

		//patch all metadata to oldTwin.
		if oldTwin.MetaData == nil {
			oldTwin.MetaData = make(map[string]*common.MetaType)
		} 
		if len(twin.MetaData) > 0 {
			for _ , meta := range twin.MetaData {
				if meta != nil {
					oldTwin.MetaData[meta.Name] = meta
				}
			}
		}

		//update desired
		if len(twin.Properties.Desired) > 0 {
			if oldTwin.Properties.Desired == nil {
				oldTwin.Properties.Desired = make(map[string]*common.TwinProperty)
			}

			for _ , prop := range twin.Properties.Desired {
				oldTwin.Properties.Desired[prop.Name] = prop
			}
		}	

		//update reported
		if len(twin.Properties.Reported) > 0 {
			if oldTwin.Properties.Reported == nil {
				oldTwin.Properties.Reported = make(map[string]*common.TwinProperty)
			}

			for _ , prop := range twin.Properties.Reported {
				oldTwin.Properties.Desired[prop.Name] = prop
			}
		}	

	} 

	return nil
}

/*
* Delete twin.
*/
func (ed *EdgeDescription) DeleteTwin(twinID string) error {
	if ed.FindTwins(twinID) != true {
		return errors.New("twin is not exists.")
	}

	ed.Lock(twinID)
	ed.Twins.Delete(twinID)
	ed.Unlock(twinID)
	ed.TwinMutex.Delete(twinID)

	return nil
} 

/*
* Get Twin
*/
func (ed *EdgeDescription) GetRawTwin(twinID string) ([]byte, error) {
	if ed.FindTwins(twinID) != true {
		return nil, errors.New("twin is not exists.")
	}

	ed.Lock(twinID)
	defer ed.Unlock(twinID)

	twin, ok := ed.getTwin(twinID) 
	if !ok {
		return nil, errors.New("twin is not exists.") 
	}

	return json.Marshal(twin)
}
