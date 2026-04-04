package eventstore

import (
	"sync"

	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"
)

// MemoryEventStore is an in-memory implementation of event store
type MemoryEventStore struct {
	events map[string][]*wallet.Event // key: wallet ID
	mu     sync.RWMutex
}

// NewMemoryEventStore creates a new in-memory event store
func NewMemoryEventStore() EventStore {
	return &MemoryEventStore{
		events: make(map[string][]*wallet.Event),
	}
}

// SaveEvent saves an event to the store
func (s *MemoryEventStore) SaveEvent(event *wallet.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.WalletID] = append(s.events[event.WalletID], event)
	return nil
}

// SaveEvents saves multiple events atomically
func (s *MemoryEventStore) SaveEvents(events []*wallet.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, event := range events {
		s.events[event.WalletID] = append(s.events[event.WalletID], event)
	}

	return nil
}

// GetEvents retrieves all events for a wallet
func (s *MemoryEventStore) GetEvents(walletID string) ([]*wallet.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[walletID]
	if !exists {
		return []*wallet.Event{}, nil
	}

	// Return copy to prevent external modifications
	result := make([]*wallet.Event, len(events))
	copy(result, events)

	return result, nil
}

// GetEventsAfterVersion retrieves events after a specific version
func (s *MemoryEventStore) GetEventsAfterVersion(walletID string, version int64) ([]*wallet.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[walletID]
	if !exists {
		return []*wallet.Event{}, nil
	}

	result := make([]*wallet.Event, 0)
	for _, event := range events {
		if event.Version > version {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetAllEvents retrieves all events in the store (for replay)
func (s *MemoryEventStore) GetAllEvents() ([]*wallet.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*wallet.Event, 0)
	for _, events := range s.events {
		result = append(result, events...)
	}

	return result, nil
}

// CountEvents returns total number of events
func (s *MemoryEventStore) CountEvents(walletID string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[walletID]
	if !exists {
		return 0, nil
	}

	return int64(len(events)), nil
}

// Clear removes all events (for testing)
func (s *MemoryEventStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(map[string][]*wallet.Event)
}
