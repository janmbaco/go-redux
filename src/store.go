package redux

import (
	"fmt"
	"github.com/janmbaco/go-infrastructure/errors"
	"github.com/janmbaco/go-infrastructure/errors/errorschecker"
	"github.com/janmbaco/go-infrastructure/eventsmanager"
	"github.com/janmbaco/go-redux/src/events"
)

type Store interface {
	GetState() interface{}
	Dispatch(Action)
	Subscribe(*func())
	Unsubscribe(*func())
	AddReducer(BusinessParam)
	RemoveReducer(string)
	GetStateOf(string) interface{}
	SubscribeTo(string, *func(interface{}))
	UnsubscribeFrom(string, *func(interface{}))
}

type store struct {
	*events.StoreSubscribeEventHandler
	errorDefer    errors.ErrorDefer
	thrower       errors.ErrorThrower
	catcher       errors.ErrorCatcher
	publisher     eventsmanager.Publisher
	reducers      map[string]Reducer
	actionsObject map[string]ActionsObject
	stateManagers map[string]*stateManager
}

func NewStore(thrower errors.ErrorThrower, catcher errors.ErrorCatcher) Store {
	errorschecker.CheckNilParameter(map[string]interface{}{"thrower": thrower, "catcher": catcher})
	subscriptions := eventsmanager.NewSubscriptions(thrower)
	return &store{
		StoreSubscribeEventHandler: events.NewStoreSubscribeEventHandler(subscriptions),
		errorDefer:                 errors.NewErrorDefer(thrower, &storeErrorPipe{}),
		reducers:                   make(map[string]Reducer),
		actionsObject:              make(map[string]ActionsObject),
		stateManagers:              make(map[string]*stateManager),
		catcher:                    catcher,
		publisher:                  eventsmanager.NewPublisher(subscriptions, catcher),
	}
}

func (s *store) AddReducer(param BusinessParam) {
	defer s.errorDefer.TryThrowError()
	errorschecker.CheckNilParameter(map[string]interface{}{"param": param})

	if _, ko := s.reducers[param.GetSelector()]; ko {
		panic(newStoreError(MultipleReducerForSelectorError, "Cannot add multiple Reducer with the same selector!"))
	}
	for _, actionsObject := range s.actionsObject {
		if actionsObject == param.GetActionsObject() {
			panic(newStoreError(MultipleReducerForActionsObjectError, "Cannot add multiple reducer with the same ActionsObject!"))
		}
	}
	s.reducers[param.GetSelector()] = param.GetReducer()
	s.actionsObject[param.GetSelector()] = param.GetActionsObject()
	if _, contains := s.stateManagers[param.GetSelector()]; !contains {
		s.stateManagers[param.GetSelector()] = newStateManager(param.GetInitialState(), param.GetSelector(), s.publisher, s.thrower, s.catcher)
	}
}

func (s *store) RemoveReducer(selector string) {
	defer s.errorDefer.TryThrowError()
	checkSelector(selector)
	s.checkStateManager(selector)
	if _, ok := s.reducers[selector]; ok {
		delete(s.reducers, selector)
	}
	if _, ok := s.actionsObject[selector]; ok {
		delete(s.actionsObject, selector)
	}
}

func (s *store) Dispatch(action Action) {
	defer s.errorDefer.TryThrowError()
	errorschecker.CheckNilParameter(map[string]interface{}{"action": action})
	var selector string
	for sel, actionObject := range s.actionsObject {
		if actionObject.Contains(action) {
			selector = sel
			break
		}
	}
	if selector == "" {
		panic(newStoreError(AnyReducerForThisActionError, "There are not any Reducers that execute this action!"))
	}

	s.stateManagers[selector].SetState((*s.reducers[selector])(s.stateManagers[selector].GetState(), action))
}

func (s *store) GetState() interface{} {
	defer s.errorDefer.TryThrowError()
	globalState := make(map[string]interface{})
	for selector, stateManager := range s.stateManagers {
		globalState[selector] = stateManager.GetState()
	}
	return globalState
}

func (s *store) GetStateOf(selector string) interface{} {
	defer s.errorDefer.TryThrowError()
	checkSelector(selector)
	s.checkStateManager(selector)
	return s.stateManagers[selector].GetState()
}

func (s *store) SubscribeTo(selector string, fn *func(interface{})) {
	defer s.errorDefer.TryThrowError()
	checkSelector(selector)
	s.checkStateManager(selector)
	s.stateManagers[selector].Subscribe(fn)
}

func (s *store) UnsubscribeFrom(selector string, fn *func(interface{})) {
	defer s.errorDefer.TryThrowError()
	checkSelector(selector)
	s.checkStateManager(selector)
	s.stateManagers[selector].UnSubscribe(fn)
}

func checkSelector(selector string) {
	if selector == "" {
		panic(newStoreError(EmptySelectorError, "The selector can not be string empty!"))
	}
}
func (s *store) checkStateManager(selector string) {
	if _, ok := s.stateManagers[selector]; !ok {
		panic(newStoreError(AnyStateBySelectorError, fmt.Sprintf("There is not any state with the selector: '%v'!", selector)))
	}
}
