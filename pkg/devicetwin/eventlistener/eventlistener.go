package eventlistener

import(
	"sync"
	"time"
	"errors"	
)

const (
	EVENT_EDGE_CREATED	= "0"
	EVENT_EDGE_ONLINE	= "1"
	EVENT_EDGE_DELETED	= "2"
	EVENT_EDGE_OFFLINE	= "3"
	EVENT_TWIN_CREATED	= "4"
	EVENT_TWIN_ONLINE	= "5"
	EVENT_TWIN_OFFLINE	= "6"
	EVENT_TWIN_UPDATE	= "7"
	EVENT_TWIN_DELETED	= "8"
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

func (el *EventListener) DeleteEventListener(){
	if el.NotifyCh != nil {
		close(el.NotifyCh)
	}
}

func (el *EventListener) MakeListenerID() string {
	return el.EdgeID+"/"+el.TwinID+"/"+el.EventId
}

func (el *EventListener) SendEventNotify() {
	event := el.EventId
	notifyCh := el.NotifyCh

	notifyCh<- event
}

func (el *EventListener) WaitEvent(callback func (interface{}))(error){
	if el.TimeOut > 0 {
		select {
		case  <-time.After(el.TimeOut):
			return errors.New("timeout!")
		case v, ok := <-el.NotifyCh:
			if !ok {
				return errors.New("channel has been closed!")
			}

			callback(v)
		}
	}else{
		v, ok := <-el.NotifyCh
		if !ok {
			return errors.New("channel has been closed!")
		}

		callback(v)
	}

	return nil
}

func (el *EventListener) WaitEventSeries(callback func (interface{}))(error){

	for{
		select {
		case  <-time.After(el.TimeOut):
			return nil
		case v, ok := <-el.NotifyCh:
			if !ok {
				return errors.New("channel has been closed!")
			}

			callback(v)
		}
	}
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
	listener.DeleteEventListener()

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

/* Match the event and dispatch it. */
func MatchEventAndDispatch(edgeID, twinID, eventID string) error {
	listenerID := edgeID+"/"+twinID+"/"+eventID

	listener := defaultListenerManager.GetEventListener(listenerID)
	if listener == nil {
		//No matched event, we just return.
		return nil
	}

	/* matched, then we dispatch the event */
	listener.SendEventNotify()

	return nil
}
