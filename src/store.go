package redux

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errorhandler"
	"github.com/janmbaco/go-infrastructure/events"
)

type Store interface {
	GetState() interface{}
	Dispatch(Action)
	Subscribe(*func())
	UnSubscribe(*func())
	AddBusinessObject(BusinessObject)
	RemoveBusinessObject(BusinessObject)
	GetStateOf(string) interface{}
	SubscribeTo(string, *func(interface{}))
	UnSubscribeFrom(string, *func(interface{}))
}

const onNewState = "onNewState"

type store struct {
	businessObjects map[string]BusinessObject
	stateManagers   map[string]*stateManager
	publisher       events.Publisher
}

func NewStore(businessObjects ...BusinessObject) Store {
	if len(businessObjects) == 0 {
		panic("There must be at least one businessObject parameter")
	}

	newStore := &store{
		businessObjects: make(map[string]BusinessObject),
		stateManagers:   make(map[string]*stateManager),
		publisher:       events.NewPublisher(),
	}
	actions := make(map[ActionsObject]bool)
	for _, bo := range businessObjects {
		if _, ko := newStore.businessObjects[bo.GetSelector()]; ko {
			panic("Cannot add multiple BusinessObject with the same StateEntity!")
		}
		if _, ko := actions[bo.GetActionsObject()]; ko {
			panic("Cannot add multiple BusinessObject with the same ActionsObject!")
		}
		newStore.businessObjects[bo.GetSelector()] = bo
		newStore.stateManagers[bo.GetSelector()] = newStateManager(newStore.publisher, bo.GetInitialState(), bo.GetSelector())
		actions[bo.GetActionsObject()] = true
	}

	return newStore
}

func (s *store) AddBusinessObject(object BusinessObject) {
	errorhandler.CheckNilParameter(map[string]interface{}{"object": object})
	if _, ko := s.businessObjects[object.GetSelector()]; ko {
		panic("Cannot add multiple BusinessObject with the same selector!")
	}
	for _, bo := range s.businessObjects {
		if bo.GetActionsObject() == object.GetActionsObject() {
			panic("Cannot add multiple BusinessObject with the same ActionsObject!")
		}
	}
	s.businessObjects[object.GetSelector()] = object
	if _, contains := s.stateManagers[object.GetSelector()]; !contains {
		s.stateManagers[object.GetSelector()] = newStateManager(s.publisher, object.GetInitialState(), object.GetSelector())
	}
}

func (s *store) RemoveBusinessObject(object BusinessObject) {
	errorhandler.CheckNilParameter(map[string]interface{}{"object": object})
	if _, ok := s.businessObjects[object.GetSelector()]; ok {
		delete(s.businessObjects, object.GetSelector())
	}
}

func (s *store) Dispatch(action Action) {
	errorhandler.CheckNilParameter(map[string]interface{}{"action": action})
	dispatched := false
	for _, bo := range s.businessObjects {
		if bo.GetActionsObject().Contains(action) {
			s.stateManagers[bo.GetSelector()].SetState(bo.GetReducer().Reduce(s.stateManagers[bo.GetSelector()].GetState(), action))
			dispatched = true
			break
		}
	}
	if !dispatched {
		panic("There is not any Reducers that execute this action!")
	}

}

func (s *store) Subscribe(subscribeFunc *func()) {
	errorhandler.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	s.publisher.Subscribe(onNewState, subscribeFunc)

}

func (s *store) UnSubscribe(subscribeFunc *func()) {
	errorhandler.CheckNilParameter(map[string]interface{}{"subscribeFunc": subscribeFunc})
	s.publisher.UnSubscribe(onNewState, subscribeFunc)
}

func (s *store) GetState() interface{} {
	var globalState interface{}
	if len(s.businessObjects) == 1 {
		for _, stateManager := range s.stateManagers {
			globalState = stateManager.GetState()
		}
	} else {
		globalState = make([]interface{}, 0)
		for _, stateManager := range s.stateManagers {
			globalState = append(globalState.([]interface{}), stateManager.GetState())
		}
	}
	return globalState
}

func (s *store) GetStateOf(selector string) interface{} {
	s.checkSelector(selector)
	return s.stateManagers[selector].GetState()
}

func (s *store) SubscribeTo(selector string, fn *func(interface{})) {
	s.checkSelector(selector)
	s.stateManagers[selector].Subscribe(fn)
}

func (s *store) UnSubscribeFrom(selector string, fn *func(interface{})) {
	s.checkSelector(selector)
	s.stateManagers[selector].UnSubscribe(fn)
}

func (s *store) checkSelector(selector string) {
	if selector == "" {
		panic("The selector can not be string empty!")
	}
	if _, ok := s.stateManagers[selector]; !ok {
		panic(fmt.Sprintf("There is not any state with the selector: '%v'!", selector))
	}
}
