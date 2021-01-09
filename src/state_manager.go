package redux

import (
	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
	"reflect"
)

type stateManager struct {
	publisher     events.Publisher
	state         reflect.Value
	typ           reflect.Type
	subscriptions map[*func(interface{})]*func()
	subscriptors  []func()
	selector      string
	isBusy        chan bool
}

func newStateManager(publisher events.Publisher, initialState interface{}, selector string) *stateManager {
	errorhandler.CheckNilParameter(map[string]interface{}{"publisher": publisher, "initialState": initialState})
	if selector == "" {
		panic("The selector can not be string empty!")
	}
	return &stateManager{publisher: publisher, state: reflect.ValueOf(initialState), typ: reflect.TypeOf(initialState), subscriptions: make(map[*func(interface{})]*func()), subscriptors: make([]func(), 0), selector: selector, isBusy: make(chan bool, 1)}

}

func (s *stateManager) GetState() interface{} {
	newState := s.state
	if s.typ.Kind() == reflect.Ptr {
		newState = reflect.New(s.typ.Elem())
		errorhandler.TryPanic(copier.Copy(newState.Interface(), s.state.Interface()))
	}
	return newState.Interface()
}

func (s *stateManager) SetState(newState interface{}) {
	s.isBusy <- true
	if !reflect.DeepEqual(newState, s.state.Interface()) {
		s.state = reflect.ValueOf(newState)
		s.publisher.Publish(onNewState)
		s.publisher.Publish(onNewState + s.selector)
	}
	<-s.isBusy
}

func (s *stateManager) Subscribe(fn *func(interface{})) {
	errorhandler.CheckNilParameter(map[string]interface{}{"fn": fn})
	s.isBusy <- true
	if _, isContained := s.subscriptions[fn]; !isContained {
		s.subscriptors = append(s.subscriptors, func() {
			(*fn)(s.GetState())
		})
		s.subscriptions[fn] = &s.subscriptors[len(s.subscriptors)-1]
		s.publisher.Subscribe(onNewState+s.selector, s.subscriptions[fn])
	}
	<-s.isBusy
}

func (s *stateManager) UnSubscribe(fn *func(interface{})) {
	errorhandler.CheckNilParameter(map[string]interface{}{"fn": fn})
	s.isBusy <- true
	order := 0
	found := false
	for function, fnEvent := range s.subscriptions {
		if function == fn {
			s.publisher.UnSubscribe(onNewState+s.selector, fnEvent)
			delete(s.subscriptions, function)
			found = true
		}
		if !found {
			order++
		}
	}
	if order < len(s.subscriptors) {
		subscriptors := make([]func(), 0)
		for i, fnEvent := range s.subscriptors {
			if i != order {
				subscriptors = append(subscriptors, fnEvent)
			}
		}
		s.subscriptors = subscriptors
	}
	<-s.isBusy
}
