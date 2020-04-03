package monitior

import(
	"sync"
	"time"
)

const (

)



type Monitior struct {
	EdgeID		string
	TwinID		string
	NotifyCh	chan interface{}
	TimeOut		time.Duration
}


func NewMonitior(edgeID, twinID	string, timeOut	time.Duration) *Monitior {
	return &Monitior{
		EdgeID: edgeID,
		TwinID: twinID,
		NotifyCh: make(chan interface{}),
		TimeOut: timeOut,	
	}
} 

/*
* Monitior manager.
*/
type MonitiorManger struct{
	MonitiorMap *sync.Map
}

var defaultMonitiorManger := &MonitiorManger{}

func NewMonitiorManger() *MonitiorManger {
	var monitiorMap sync.Map
	return &MonitiorManger{
		MonitiorMap: &monitiorMap,
	}
} 

func AddMonitior(edgeID, twinID	string, timeOut	time.Duration)(string, error){

} 

