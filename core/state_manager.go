package redux

import (
	"github.com/janmbaco/copier"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
	"reflect"
)

const onNewStae = "onNewState"

type StateManager interface {
	GetState() interface{}
	SetState(interface{})
	Subscribe(fn func())
}

type stateManager struct {
	publisher events.EventPublisher
	state     reflect.Value
	typ       reflect.Type
}

func NewStateManager(publisher events.EventPublisher, stateEntity StateEntity) StateManager {
	errorhandler.CheckNilParameter(map[string]interface{}{"publisher": publisher, "stateEntity": stateEntity})
	return &stateManager{publisher: publisher, state: reflect.ValueOf(stateEntity.GetInitialState()), typ: reflect.TypeOf(stateEntity.GetInitialState())}

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
	s.state = reflect.ValueOf(newState)
	s.publisher.Publish(onNewStae)
}

func (s *stateManager) Subscribe(fn func()) {
	s.publisher.Subscribe(onNewStae, fn)
}
