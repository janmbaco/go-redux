package snapshots

import (
	"sync"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"
)

// SnapshotStore stores wallet snapshots
type SnapshotStore struct {
	snapshots map[string]*wallet.SnapshotData // key: wallet ID
	mu        sync.RWMutex
}

// NewSnapshotStore creates a new snapshot store
func NewSnapshotStore() *SnapshotStore {
	return &SnapshotStore{
		snapshots: make(map[string]*wallet.SnapshotData),
	}
}

// SaveSnapshot saves a snapshot
func (s *SnapshotStore) SaveSnapshot(walletID string, snapshot *wallet.SnapshotData) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.snapshots[walletID] = snapshot
}

// GetSnapshot retrieves a snapshot
func (s *SnapshotStore) GetSnapshot(walletID string) (*wallet.SnapshotData, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot, exists := s.snapshots[walletID]
	return snapshot, exists
}

// Snapshotter handles automatic snapshot creation
type Snapshotter struct {
	store          *SnapshotStore
	eventThreshold int64 // Create snapshot every N events
	logger         logs.Logger
}

// NewSnapshotter creates a new snapshotter
func NewSnapshotter(store *SnapshotStore, eventThreshold int64, logger logs.Logger) *Snapshotter {
	return &Snapshotter{
		store:          store,
		eventThreshold: eventThreshold,
		logger:         logger,
	}
}

// CheckAndCreateSnapshot checks if snapshot is needed and creates it
func (s *Snapshotter) CheckAndCreateSnapshot(w *wallet.Wallet, eventCount int64, lastEventID string) {
	if eventCount%s.eventThreshold == 0 {
		s.logger.Info("[Snapshotter] Creating snapshot for wallet " + w.ID)

		snapshot := &wallet.SnapshotData{
			Balance:     w.Balance,
			Version:     w.Version,
			LastEventID: lastEventID,
		}

		s.store.SaveSnapshot(w.ID, snapshot)
	}
}

// LoadFromSnapshot loads wallet state from snapshot and applies remaining events
func (s *Snapshotter) LoadFromSnapshot(walletID string, allEvents []*wallet.Event) (*wallet.Wallet, error) {
	snapshot, exists := s.store.GetSnapshot(walletID)

	if !exists {
		// No snapshot, replay all events
		s.logger.Info("[Snapshotter] No snapshot found for " + walletID + ", replaying all events")
		return s.replayAllEvents(allEvents)
	}

	s.logger.Info("[Snapshotter] Loading from snapshot for " + walletID)

	// Create wallet from snapshot
	w := &wallet.Wallet{
		ID:        walletID,
		Balance:   snapshot.Balance,
		Version:   snapshot.Version,
		UpdatedAt: time.Now(),
	}

	// Apply events after snapshot version
	for _, event := range allEvents {
		if event.Version > snapshot.Version {
			w.ApplyEvent(event)
		}
	}

	return w, nil
}

func (s *Snapshotter) replayAllEvents(events []*wallet.Event) (*wallet.Wallet, error) {
	if len(events) == 0 {
		return nil, wallet.ErrWalletNotFound
	}

	// First event should be WalletCreated
	firstEvent := events[0]
	if firstEvent.Type != wallet.EventWalletCreated {
		return nil, wallet.ErrWalletNotFound
	}

	data := firstEvent.Data
	w := &wallet.Wallet{
		ID:        firstEvent.WalletID,
		OwnerID:   data["owner_id"].(string),
		Currency:  data["currency"].(string),
		Balance:   data["initial_balance"].(float64),
		Version:   firstEvent.Version,
		CreatedAt: firstEvent.Timestamp,
		UpdatedAt: firstEvent.Timestamp,
	}

	// Apply remaining events
	for i := 1; i < len(events); i++ {
		w.ApplyEvent(events[i])
	}

	return w, nil
}
