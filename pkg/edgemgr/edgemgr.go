/*
* The edgemgr has some function as below:
* 1. Detect  the edge cluster on remote.
* 2. Keepalive between cloud and edge sides.
* 3. Manage the edge on remote.
*/
package edgemgr


import (

)

type EdgeDescription struct {
	ID		string
	Name 	string
	Description	string
	State	string
}

type EdgeManager struct {
	EdgeCache map[string]EdgeDescription	
}

func NewEdgeManager() *EdgeManager {
	return &EdgeManager{
		EdgeCache: make(map[string]EdgeDescription),
	}
}

func (em *EdgeManager) CreateEdge(id string) {

}
