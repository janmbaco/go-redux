package redux

import (
	"fmt"
	"sync"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2/events"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// Store represents the Redux store that holds application state.
// The store manages state through modules and notifies subscribers of changes.
//
// Example:
//
//	store := redux.NewStore[map[string]any](make(map[string]any), logger)
//	store.AddModule(counterModule)
//	store.Dispatch(IncrementAction.With(5))
type Store[S any] interface {
	// GetState returns the current state
	GetState() S

	// Dispatch sends an action to the store
	Dispatch(action any) error

	// AddModule adds a module to the store
	AddModule(module any) error

	// Subscribe registers a callback for state changes
	Subscribe(callback *func(S)) error

	// Unsubscribe removes a callback
	Unsubscribe(callback *func(S)) error

	// Close shuts down the store
	Close()
}

// store is the concrete implementation of Store
type store[S any] struct {
	mu            sync.RWMutex
	state         S
	modules       []any
	stateEventMgr *events.StateEventManager[S]
	closed        bool
	logger        logs.Logger
}

// NewStore creates a new Store with the given initial state
func NewStore[S any](initialState S, logger logs.Logger) Store[S] {
	return &store[S]{
		state:         initialState,
		modules:       make([]any, 0),
		stateEventMgr: events.NewStateEventManager[S](logger),
		closed:        false,
		logger:        logger,
	}
}

// GetState returns the current state
func (s *store[S]) GetState() S {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// Dispatch sends an action to the store.
// The action is processed by all modules that can handle it.
func (s *store[S]) Dispatch(action any) error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return fmt.Errorf("store is closed")
	}

	// Handle map[string]any state specially (common case)
	stateMap, isMap := any(s.state).(map[string]any)

	var err error
	stateChanged := false

	// Process action through all modules
	for _, mod := range s.modules {
		// Use reflection to call CanHandle and Handle methods
		canHandle, handleErr := handlers.InvokeCanHandle(mod, action)
		if handleErr != nil {
			continue
		}

		if !canHandle {
			continue
		}

		// Get module selector
		selector := handlers.InvokeGetSelector(mod)

		if isMap && selector != "" {
			// Handle map state with selector
			moduleState := stateMap[selector]
			newModuleState, handleErr := handlers.InvokeHandle(mod, moduleState, action)
			if handleErr != nil {
				err = handleErr
				continue
			}
			stateMap[selector] = newModuleState
			stateChanged = true
		} else {
			// Handle non-map state
			newState, handleErr := handlers.InvokeHandle(mod, s.state, action)
			if handleErr != nil {
				err = handleErr
				continue
			}
			// Type assertion to S
			if typedState, ok := newState.(S); ok {
				s.state = typedState
				stateChanged = true
			}
		}
	}

	currentState := s.state
	s.mu.Unlock()

	// Notify subscribers if state changed
	if stateChanged {
		s.stateEventMgr.Publish(currentState)
	}

	return err
}

// AddModule adds a handler to the store.
// For map[string]any states, the handler's initial state is added under its selector key.
func (s *store[S]) AddModule(handler any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("store is closed")
	}

	// Initialize handler state if using map
	stateMap, isMap := any(s.state).(map[string]any)
	if isMap {
		selector := handlers.InvokeGetSelector(handler)
		if selector != "" {
			initialState := handlers.InvokeGetInitialState(handler)
			stateMap[selector] = initialState
		}
	}

	s.modules = append(s.modules, handler)
	return nil
}

// Subscribe registers a callback for state changes
func (s *store[S]) Subscribe(callback *func(S)) error {
	return s.stateEventMgr.Subscribe(callback)
}

// Unsubscribe removes a callback
func (s *store[S]) Unsubscribe(callback *func(S)) error {
	return s.stateEventMgr.Unsubscribe(callback)
}

// Close shuts down the store
func (s *store[S]) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
}
