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
	Subscribe(fn SubscribeFunc)
	UnSubscribe(fn SubscribeFunc)
}

type stateManager struct {
	publisher     events.Publisher
	state         reflect.Value
	typ           reflect.Type
	subscriptions map[uintptr]func()
	isBusy        chan bool
}

func NewStateManager(publisher events.Publisher, stateEntity StateEntity) StateManager {
	errorhandler.CheckNilParameter(map[string]interface{}{"publisher": publisher, "stateEntity": stateEntity})
	return &stateManager{publisher: publisher, state: reflect.ValueOf(stateEntity.GetInitialState()), typ: reflect.TypeOf(stateEntity.GetInitialState()), subscriptions: make(map[uintptr]func()), isBusy: make(chan bool, 1)}

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

func (s *stateManager) Subscribe(fn SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"fn": fn})
	s.isBusy <- true
	pointer := reflect.ValueOf(fn).Pointer()
	if _, isContained := s.subscriptions[pointer]; !isContained {
		fnEvent := func() {
			fn(s.GetState())
		}
		s.subscriptions[pointer] = fnEvent
		s.publisher.Subscribe(onNewStae, fnEvent)
	}
	<-s.isBusy
}

func (s *stateManager) UnSubscribe(fn SubscribeFunc) {
	errorhandler.CheckNilParameter(map[string]interface{}{"fn": fn})
	s.isBusy <- true
	subscriptions := make(map[uintptr]func())
	for pointer, fnEvent := range s.subscriptions {
		if pointer != reflect.ValueOf(fn).Pointer() {
			subscriptions[pointer] = fnEvent
		} else {
			s.publisher.UnSubscribe(onNewStae, fnEvent)
		}
	}
	s.subscriptions = subscriptions
	<-s.isBusy
}
