package eventlistener

import(
	"sync"
	"time"
	"errors"	
)

const (
	EVENT_	= 1	
)

type EventListener struct {
	EdgeID		string
	TwinID		string
	EventId		string	
	NotifyCh	chan interface{}
	TimeOut		time.Duration
}

func NewEventListener(edgeID, twinID, eventID string, 
			timeOut time.Duration) *EventListener {
		
	return &EventListener{
		EdgeID: edgeID,
		TwinID:	twinID, 
		EventId: eventID,	
		NotifyCh: make(chan interface{}),
		TimeOut: timeOut,	
	}	
}

func (el *EventListener) MakeListenerID() string {
	return el.EdgeID+"/"+el.TwinID+"/"+el.EventId
}

/*
* Event Listener manager
*/
type EventListenerManager struct {
	EventListenerMap *sync.Map
}

var defaultListenerManager = &EventListenerManager{}


func init() {
	var listenerMap sync.Map
	defaultListenerManager.EventListenerMap = &listenerMap	
}

func (elm *EventListenerManager) GetEventListener(listenerID string) *EventListener {
	v, exist := elm.EventListenerMap.Load(listenerID)
	if !exist {
		return nil	
	}
	
	eventListener, isThisType := v.(*EventListener)
	if !isThisType {
		return nil
	}			
	
	return eventListener
}

func (elm *EventListenerManager) PutEventListener(listener *EventListener) error {
	if listener == nil {
		return errors.New("listener is nil")
	}
	
	listenerID := listener.MakeListenerID()
	if 	listenerID == "//" {
		return errors.New("invalid listener id")
	}
	
	if elm.GetEventListener(listenerID) != nil 	{
		return errors.New("listener has exists")
	}	
		
	elm.EventListenerMap.Store(listenerID, listener)
	
	return nil
}

func (elm *EventListenerManager) DeleteEventListener(listener *EventListener) error {
	if listener == nil {
		return errors.New("listener is nil")
	}
	
	listenerID := listener.MakeListenerID()
	if 	listenerID == "//" {
		return errors.New("invalid listener id")
	}
	
	if elm.GetEventListener(listenerID) == nil 	{
		return errors.New("listener not exists")
	}
	
	elm.EventListenerMap.Delete(listenerID)
	
	return nil	
}

/* register the event listener. */
func RegisterEventListener(edgeID, twinID, eventID string, 
				timeOut time.Duration) (*EventListener, error) {

	eventListener := NewEventListener(edgeID, twinID, 
			eventID, timeOut)
	
	/* Add event listener into defaultListenerManager */
	err := defaultListenerManager.PutEventListener(eventListener)
	if err != nil {
		return nil, err
	}
		
	return 	eventListener, nil	
}

/* unregister the event listener.*/
func UnregisterEventListener(listener **EventListener) error {
	err := defaultListenerManager.DeleteEventListener(*listener)
	*listener = nil
	
	return err
}
