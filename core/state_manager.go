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
	state     interface{}
	typ       reflect.Type
}

func NewStateManager(publisher events.EventPublisher, stateEntity *StateEntity) StateManager {
	errorhandler.CheckNilParameter(map[string]interface{}{"publisher": publisher, "stateEntity": stateEntity})
	return &stateManager{publisher: publisher, state: stateEntity.InitialState, typ: reflect.TypeOf(stateEntity.InitialState)}
}

func (s *stateManager) GetState() interface{} {
	newState := s.state
	if s.typ.Kind() == reflect.Ptr {
		newState := reflect.New(s.typ).Interface()
		errorhandler.TryPanic(copier.Copy(newState, s.state))
	}
	return newState
}

func (s *stateManager) SetState(newState interface{}) {
	s.state = newState
	s.publisher.Publish(onNewStae)
}

func (s *stateManager) Subscribe(fn func()) {
	s.publisher.Subscribe(onNewStae, fn)
}
