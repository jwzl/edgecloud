/*
* The edgemgr has some function as below:
* 1. Detect  the edge cluster on remote.
* 2. Keepalive between cloud and edge sides.
* 3. Manage the edge on remote.
*/
package edge


import (
	"sync"
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
	twins		sync.Map
	twinMutex	sync.Map

	return &EdgeDescription{
		ID: edgeID,
		State: EdgeStateInitial,
		deviceIDs: make([]string, 0),
		Twins: &twins,
		TwinMutex: twinMutex, 
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


func (ed *EdgeDescription) UpdateTwin(twin common.DigitalTwin) error {
	twinID := twin.ID
	if ed.FindTwins(twinID) != true {
		return errors.New("No such this twin")
	}

	_ , exist := ed.TwinMutex.Load(twinID)
	if !exist {
		//Create twin
		
	} else {
		// update twin
		ed.Lock(twinID)
		defer ed.Unlock(twinID)	

	} 
}
