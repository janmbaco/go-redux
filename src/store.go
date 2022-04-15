package redux

import (
	"fmt"
	"reflect"

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
	publisher     eventsmanager.Publisher
	reducers      map[string]Reducer
	actionsObject map[string]ActionsObject
	stateManagements map[string]StateManagement
	stateManagementFactory StateManagementFactory
}

func NewStore(errorDefer errors.ErrorDefer, subscriptions eventsmanager.Subscriptions, publisher eventsmanager.Publisher, stateManagementFactory StateManagementFactory) Store {
	errorschecker.CheckNilParameter(map[string]interface{}{"errorDefer": errorDefer, "subscriptions": subscriptions, "publisher": publisher, "stateManagementResolver":stateManagementFactory})
	return &store{
		StoreSubscribeEventHandler: events.NewStoreSubscribeEventHandler(subscriptions),
		errorDefer:                 errorDefer,
		reducers:                   make(map[string]Reducer),
		actionsObject:              make(map[string]ActionsObject),
		stateManagements:           make(map[string]StateManagement),
		publisher:                  publisher,
		stateManagementFactory:    stateManagementFactory,
	}
}

func (s *store) AddReducer(param BusinessParam) {
	defer s.errorDefer.TryThrowError(s.errorPipe)
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
	if _, contains := s.stateManagements[param.GetSelector()]; !contains {
		s.stateManagements[param.GetSelector()] = s.stateManagementFactory.Create(StateManagementFactoryParamter{param.GetInitialState(), param.GetSelector(), s.publisher})
	}
}

func (s *store) RemoveReducer(selector string) {
	defer s.errorDefer.TryThrowError(s.errorPipe)
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
	defer s.errorDefer.TryThrowError(s.errorPipe)
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

	s.stateManagements[selector].SetState((*s.reducers[selector])(s.stateManagements[selector].GetState(), action))
}

func (s *store) GetState() interface{} {
	defer s.errorDefer.TryThrowError(s.errorPipe)
	globalState := make(map[string]interface{})
	for selector, stateManager := range s.stateManagements {
		globalState[selector] = stateManager.GetState()
	}
	return globalState
}

func (s *store) GetStateOf(selector string) interface{} {
	defer s.errorDefer.TryThrowError(s.errorPipe)
	checkSelector(selector)
	s.checkStateManager(selector)
	return s.stateManagements[selector].GetState()
}

func (s *store) SubscribeTo(selector string, fn *func(interface{})) {
	defer s.errorDefer.TryThrowError(s.errorPipe)
	checkSelector(selector)
	s.checkStateManager(selector)
	s.stateManagements[selector].Subscribe(fn)
}

func (s *store) UnsubscribeFrom(selector string, fn *func(interface{})) {
	defer s.errorDefer.TryThrowError(s.errorPipe)
	checkSelector(selector)
	s.checkStateManager(selector)
	s.stateManagements[selector].UnSubscribe(fn)
}

func checkSelector(selector string) {
	if selector == "" {
		panic(newStoreError(EmptySelectorError, "The selector can not be string empty!"))
	}
}
func (s *store) checkStateManager(selector string) {
	if _, ok := s.stateManagements[selector]; !ok {
		panic(newStoreError(AnyStateBySelectorError, fmt.Sprintf("There is not any state with the selector: '%v'!", selector)))
	}
}

func (s *store) errorPipe(err error) error {
	resultError := err

	if errType := reflect.Indirect(reflect.ValueOf(err)).Type(); !errType.Implements(reflect.TypeOf((*StoreError)(nil)).Elem()) {
		errorType := UnexpectedStoreError
		resultError = &storeError{
			CustomizableError: errors.CustomizableError{
				Message:       err.Error(),
				InternalError: err,
			},
			ErrorType: errorType,
		}
	}
	return resultError
}
