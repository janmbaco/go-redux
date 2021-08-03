package redux

import (
	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
	"github.com/janmbaco/go-redux/src/events"
	"reflect"
)

type stateManager struct {
	*events.SelectorSubscribeEventHandler
	storePublisher    eventsmanager.Publisher
	selectorPublisher eventsmanager.Publisher
	state             reflect.Value
	typ               reflect.Type
	subscriptors      []func()
	selector          string
}

func newStateManager(initialState interface{}, selector string, storePublisher eventsmanager.Publisher, thrower errors.ErrorThrower, catcher errors.ErrorCatcher) *stateManager {
	errorschecker.CheckNilParameter(map[string]interface{}{"publisher": storePublisher, "initialState": initialState})
	if selector == "" {
		panic("The selector can not be string empty!")
	}
	subscriptions := eventsmanager.NewSubscriptions(thrower)
	return &stateManager{
		SelectorSubscribeEventHandler: events.NewSelectorSubscribeEventHandler(subscriptions),
		storePublisher:                storePublisher,
		state:                         reflect.ValueOf(initialState),
		typ:                           reflect.TypeOf(initialState),
		selectorPublisher:             eventsmanager.NewPublisher(subscriptions, catcher),
		selector:                      selector,
	}

}

func (s *stateManager) GetState() interface{} {
	newState := s.state
	if s.typ.Kind() == reflect.Ptr {
		newState = reflect.New(s.typ.Elem())
		errorschecker.TryPanic(copier.Copy(newState.Interface(), s.state.Interface()))
	}
	return newState.Interface()
}

func (s *stateManager) SetState(newState interface{}) {
	if !reflect.DeepEqual(newState, s.state.Interface()) {
		s.state = reflect.ValueOf(newState)
		s.storePublisher.Publish(&events.StoreSubscribeEvent{})
		s.selectorPublisher.Publish(&events.SelectorSubscribeEvent{State: s.GetState()})
	}
}
