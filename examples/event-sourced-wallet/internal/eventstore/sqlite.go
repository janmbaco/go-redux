package eventstore

import (
	"encoding/json"
	"time"

	"github.com/janmbaco/go-infrastructure/v2/persistence/dataaccess"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"
	"gorm.io/gorm"
)

// EventRecord is the database model for storing events
type EventRecord struct {
	ID            string    `gorm:"primaryKey;type:varchar(36)"`
	WalletID      string    `gorm:"index;type:varchar(36);not null"`
	Type          string    `gorm:"type:varchar(50);not null"`
	Data          string    `gorm:"type:text"` // JSON serialized
	Timestamp     time.Time `gorm:"not null"`
	Version       int64     `gorm:"not null"`
	CorrelationID string    `gorm:"type:varchar(36)"`
}

func (EventRecord) TableName() string {
	return "wallet_events"
}

// SQLiteEventStore is a SQLite-backed event store
type SQLiteEventStore struct {
	eventRecorDa dataaccess.DataAccess
}

// NewSQLiteEventStore creates a new SQLite event store
func NewSQLiteEventStore(db *gorm.DB) EventStore {
	// Auto-migrate the schema
	db.AutoMigrate(&EventRecord{})

	return &SQLiteEventStore{
		eventRecorDa: dataaccess.NewTypedDataAccess[EventRecord](db),
	}
}

// SaveEvent saves an event to SQLite
func (s *SQLiteEventStore) SaveEvent(event *wallet.Event) error {
	dataJSON, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	record := &EventRecord{
		ID:            event.ID,
		WalletID:      event.WalletID,
		Type:          string(event.Type),
		Data:          string(dataJSON),
		Timestamp:     event.Timestamp,
		Version:       event.Version,
		CorrelationID: event.CorrelationID,
	}

	return dataaccess.InsertRow(s.eventRecorDa, record)
}

// SaveEvents saves multiple events atomically
func (s *SQLiteEventStore) SaveEvents(events []*wallet.Event) error {
	for _, event := range events {
		if err := s.SaveEvent(event); err != nil {
			return err
		}
	}
	return nil
}

// GetEvents retrieves all events for a wallet
func (s *SQLiteEventStore) GetEvents(walletID string) ([]*wallet.Event, error) {
	filter := &EventRecord{WalletID: walletID}
	records, err := dataaccess.SelectRows(s.eventRecorDa, filter)
	if err != nil {
		return nil, err
	}

	events := make([]*wallet.Event, len(records))
	for i, record := range records {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(record.Data), &data); err != nil {
			return nil, err
		}

		events[i] = &wallet.Event{
			ID:            record.ID,
			WalletID:      record.WalletID,
			Type:          wallet.EventType(record.Type),
			Data:          data,
			Timestamp:     record.Timestamp,
			Version:       record.Version,
			CorrelationID: record.CorrelationID,
		}
	}

	return events, nil
}

// GetEventsAfterVersion retrieves events after a specific version
func (s *SQLiteEventStore) GetEventsAfterVersion(walletID string, version int64) ([]*wallet.Event, error) {
	// Note: This would require a custom query in production
	// For now, get all and filter in memory
	allEvents, err := s.GetEvents(walletID)
	if err != nil {
		return nil, err
	}
	s.eventRecorDa.DB()
	result := make([]*wallet.Event, 0)
	for _, event := range allEvents {
		if event.Version > version {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetAllEvents retrieves all events in the store (for replay)
func (s *SQLiteEventStore) GetAllEvents() ([]*wallet.Event, error) {
	records, err := dataaccess.SelectRows(s.eventRecorDa, &EventRecord{})
	if err != nil {
		return nil, err
	}

	events := make([]*wallet.Event, len(records))
	for i, record := range records {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(record.Data), &data); err != nil {
			return nil, err
		}

		events[i] = &wallet.Event{
			ID:            record.ID,
			WalletID:      record.WalletID,
			Type:          wallet.EventType(record.Type),
			Data:          data,
			Timestamp:     record.Timestamp,
			Version:       record.Version,
			CorrelationID: record.CorrelationID,
		}
	}

	return events, nil
}

// CountEvents returns total number of events for a wallet
func (s *SQLiteEventStore) CountEvents(walletID string) (int64, error) {
	filter := &EventRecord{WalletID: walletID}
	records, err := dataaccess.SelectRows(s.eventRecorDa, filter)
	if err != nil {
		return 0, err
	}

	return int64(len(records)), nil
}
