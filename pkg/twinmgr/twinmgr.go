package twinmgr

import (

)

type TwinManager struct {
	//edgeID: DeviceID  
	edgeTwins map[string]string
}

func NewTwinManager() *TwinManager {
	return &TwinManager{
		edgeTwins: make(map[string]string, 0),
	}
} 

func(tm *TwinManager) Start(){

	for {
	
		select {
		
	
		}	

	}
}


