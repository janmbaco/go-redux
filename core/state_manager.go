package redux

import (
	"reflect"

	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
)

type SubscribeFunc func(newState interface{})

const onNewStae = "onNewState"

type StateManager interface {
	GetState() interface{}
	SetState(interface{})
	Subscribe(fn *SubscribeFunc)
	UnSubscribe(fn *SubscribeFunc)
}

type stateManager struct {
	publisher     events.Publisher
	state         reflect.Value
	typ           reflect.Type
	subscriptions map[*SubscribeFunc]*func()
	subscriptors  []func()
	isBusy        chan bool
}

func NewStateManager(publisher events.Publisher, stateEntity StateEntity) StateManager {
	errorhandler.CheckNilParameter(map[string]interface{}{"publisher": publisher, "stateEntity": stateEntity})
	return &stateManager{publisher: publisher, state: reflect.ValueOf(stateEntity.GetInitialState()), typ: reflect.TypeOf(stateEntity.GetInitialState()), subscriptions: make(map[*SubscribeFunc]*func()), subscriptors: make([]func(), 0), isBusy: make(chan bool, 1)}

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
	s.state = reflect.ValueOf(newState)
	s.publisher.Publish(onNewStae)
	<-s.isBusy
}

func (s *stateManager) Subscribe(fn *SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"fn": fn})
	s.isBusy <- true
	if _, isContained := s.subscriptions[fn]; !isContained {
		s.subscriptors = append(s.subscriptors, func() {
			(*fn)(s.GetState())
		})
		s.subscriptions[fn] = &s.subscriptors[len(s.subscriptors)-1]
		s.publisher.Subscribe(onNewStae, s.subscriptions[fn])
	}
	<-s.isBusy
}

func (s *stateManager) UnSubscribe(fn *SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"fn": fn})
	s.isBusy <- true
	subscriptions := make(map[*SubscribeFunc]*func())
	order := 0
	found := false
	for funtion, fnEvent := range s.subscriptions {
		if funtion != fn {
			subscriptions[funtion] = fnEvent
			if !found {
				order++
			}
		} else {
			s.publisher.UnSubscribe(onNewStae, fnEvent)
			found = true
		}
	}
	s.subscriptions = subscriptions

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
