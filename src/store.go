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
	AddReducer(BusinessParam)
	RemoveReducer(string)
	GetStateOf(string) interface{}
	SubscribeTo(string, *func(interface{}))
	UnSubscribeFrom(string, *func(interface{}))
}

const onNewState = "onNewState"

type store struct {
	reducers      map[string]Reducer
	actionsObject map[string]ActionsObject
	stateManagers map[string]*stateManager
	publisher     events.Publisher
}

func NewStore(params ...BusinessParam) Store {
	if len(params) == 0 {
		panic("There must be at least one businessParam parameter")
	}

	newStore := &store{
		reducers:      make(map[string]Reducer),
		actionsObject: make(map[string]ActionsObject),
		stateManagers: make(map[string]*stateManager),
		publisher:     events.NewPublisher(),
	}
	for _, param := range params {
		if _, ko := newStore.actionsObject[param.GetSelector()]; ko {
			panic("Cannot add multiple Reducer with the same ActionsObject!")
		}
		if _, ko := newStore.reducers[param.GetSelector()]; ko {
			panic("Cannot add multiple Reducer with the same Selector!")
		}
		newStore.actionsObject[param.GetSelector()] = param.GetActionsObject()
		newStore.reducers[param.GetSelector()] = param.GetReducer()
		newStore.stateManagers[param.GetSelector()] = newStateManager(newStore.publisher, param.GetInitialState(), param.GetSelector())
	}

	return newStore
}

func (s *store) AddReducer(param BusinessParam) {
	errorhandler.CheckNilParameter(map[string]interface{}{"param": param})
	if _, ko := s.reducers[param.GetSelector()]; ko {
		panic("Cannot add multiple Reducer with the same selector!")
	}
	for _, actionsObject := range s.actionsObject {
		if actionsObject == param.GetActionsObject() {
			panic("Cannot add multiple reducer with the same ActionsObject!")
		}
	}
	s.reducers[param.GetSelector()] = param.GetReducer()
	s.actionsObject[param.GetSelector()] = param.GetActionsObject()
	if _, contains := s.stateManagers[param.GetSelector()]; !contains {
		s.stateManagers[param.GetSelector()] = newStateManager(s.publisher, param.GetInitialState(), param.GetSelector())
	}
}

func (s *store) RemoveReducer(selector string) {
	checkSelector(selector)
	s.checkReducer(selector)
	if _, ok := s.reducers[selector]; ok {
		delete(s.reducers, selector)
	}
	if _, ok := s.actionsObject[selector]; ok {
		delete(s.actionsObject, selector)
	}
}

func (s *store) Dispatch(action Action) {
	errorhandler.CheckNilParameter(map[string]interface{}{"action": action})
	var selector string
	for sel, actionObject := range s.actionsObject {
		if actionObject.Contains(action) {
			selector = sel
			break
		}
	}
	if selector == "" {
		panic("There are not any Reducers that execute this action!")
	}

	s.stateManagers[selector].SetState(s.reducers[selector](s.stateManagers[selector].GetState(), action))
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
	globalState := make(map[string]interface{})
	for selector, stateManager := range s.stateManagers {
		globalState[selector] = stateManager.GetState()
	}
	return globalState
}

func (s *store) GetStateOf(selector string) interface{} {
	checkSelector(selector)
	s.checkStateManager(selector)
	return s.stateManagers[selector].GetState()
}

func (s *store) SubscribeTo(selector string, fn *func(interface{})) {
	checkSelector(selector)
	s.checkStateManager(selector)
	s.stateManagers[selector].Subscribe(fn)
}

func (s *store) UnSubscribeFrom(selector string, fn *func(interface{})) {
	checkSelector(selector)
	s.checkStateManager(selector)
	s.stateManagers[selector].UnSubscribe(fn)
}

func checkSelector(selector string) {
	if selector == "" {
		panic("The selector can not be string empty!")
	}
}
func (s *store) checkStateManager(selector string) {
	if _, ok := s.stateManagers[selector]; !ok {
		panic(fmt.Sprintf("There is not any state with the selector: '%v'!", selector))
	}
}

func (s *store) checkReducer(selector string) {
	if _, ok := s.stateManagers[selector]; !ok {
		panic(fmt.Sprintf("There is not any state with the selector: '%v'!", selector))
	}
}
