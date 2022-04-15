package redux

import (
	"reflect"

	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
	"github.com/janmbaco/go-redux/src/events"
	"github.com/jinzhu/copier"
)
type StateManagement interface{
	Subscribe(subscription *func(state interface{}))
	UnSubscribe(subscription *func(state interface{}))
	GetState() interface{}
	SetState(newState interface{})
}

type stateManagement struct {
	*events.SelectorSubscribeEventHandler
	storePublisher    eventsmanager.Publisher
	selectorPublisher eventsmanager.Publisher
	state             reflect.Value
	typ               reflect.Type
	subscriptors      []func()
	selector          string
}

func NewStateManager(initialState interface{}, selector string, storePublisher eventsmanager.Publisher, subscriptions eventsmanager.Subscriptions, selectorPublisher eventsmanager.Publisher) StateManagement {
	errorschecker.CheckNilParameter(map[string]interface{}{"initialState": initialState, "selector": selector, "storePublisher": storePublisher, "subscriptions": subscriptions, "selectorPublisher": selectorPublisher})
	if selector == "" {
		panic("The selector can not be string empty!")
	}
	return &stateManagement{
		SelectorSubscribeEventHandler: events.NewSelectorSubscribeEventHandler(subscriptions),
		storePublisher:                storePublisher,
		state:                         reflect.ValueOf(initialState),
		typ:                           reflect.TypeOf(initialState),
		selectorPublisher:             selectorPublisher,
		selector:                      selector,
	}

}

func (s *stateManagement) GetState() interface{} {
	newState := s.state
	if s.typ.Kind() == reflect.Ptr {
		newState = reflect.New(s.typ.Elem())
		errorschecker.TryPanic(copier.Copy(newState.Interface(), s.state.Interface()))
	}
	return newState.Interface()
}

func (s *stateManagement) SetState(newState interface{}) {
	if !reflect.DeepEqual(newState, s.state.Interface()) {
		s.state = reflect.ValueOf(newState)
		s.storePublisher.Publish(&events.StoreSubscribeEvent{})
		s.selectorPublisher.Publish(&events.SelectorSubscribeEvent{State: s.GetState()})
	}
}
