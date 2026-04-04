package events

import (
	"reflect"
	"sync"

	"github.com/janmbaco/go-infrastructure/v2/eventsmanager"
	"github.com/janmbaco/go-infrastructure/v2/logs"
)

// StateChangedEvent represents a state change in the store
type StateChangedEvent[S any] struct {
	state S
}

// GetEventArgs returns the event args
func (e StateChangedEvent[S]) GetEventArgs() StateChangedEvent[S] {
	return e
}

// StopPropagation indicates if propagation should stop
func (e StateChangedEvent[S]) StopPropagation() bool {
	return false
}

// IsParallelPropagation indicates if should publish in parallel
func (e StateChangedEvent[S]) IsParallelPropagation() bool {
	return true
}

// GetState returns the state from the event
func (e StateChangedEvent[S]) GetState() S {
	return e.state
}

// StateEventManager manages state change events for a store
type StateEventManager[S any] struct {
	eventManager *eventsmanager.EventManager
	publisher    eventsmanager.Publisher[StateChangedEvent[S]]
	mu           sync.RWMutex
	handlers     map[uintptr]func(S)
}

// NewStateEventManager creates a new state event manager
func NewStateEventManager[S any](logger logs.Logger) *StateEventManager[S] {
	em := eventsmanager.NewEventManager()
	subs := eventsmanager.NewSubscriptions[StateChangedEvent[S]]()
	pub := eventsmanager.NewPublisher(subs, logger)
	eventsmanager.Register(em, pub)

	sem := &StateEventManager[S]{
		eventManager: em,
		publisher:    pub,
		handlers:     make(map[uintptr]func(S)),
	}

	// Register single wrapper that dispatches to all handlers
	subs.Add(func(event StateChangedEvent[S]) {
		for _, handler := range sem.snapshotHandlers() {
			handler(event.GetState())
		}
	})

	return sem
}

// Subscribe subscribes a handler to state changes
func (sem *StateEventManager[S]) Subscribe(handler *func(S)) error {
	pointer := reflect.ValueOf(handler).Pointer()

	sem.mu.Lock()
	defer sem.mu.Unlock()

	if _, exists := sem.handlers[pointer]; exists {
		return nil
	}

	// Store the dereferenced handler directly
	sem.handlers[pointer] = *handler

	return nil
}

// Unsubscribe removes a handler from state changes
func (sem *StateEventManager[S]) Unsubscribe(handler *func(S)) error {
	pointer := reflect.ValueOf(handler).Pointer()
	sem.mu.Lock()
	delete(sem.handlers, pointer)
	sem.mu.Unlock()
	return nil
}

// Publish publishes a state change event
func (sem *StateEventManager[S]) Publish(state S) {
	event := StateChangedEvent[S]{state: state}
	eventsmanager.Publish(sem.eventManager, event)
}

func (sem *StateEventManager[S]) snapshotHandlers() []func(S) {
	sem.mu.RLock()
	defer sem.mu.RUnlock()

	snapshot := make([]func(S), 0, len(sem.handlers))
	for _, handler := range sem.handlers {
		snapshot = append(snapshot, handler)
	}

	return snapshot
}
