package eventstore

import "github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"

// EventStore defines the interface for event storage
type EventStore interface {
	SaveEvent(event *wallet.Event) error
	SaveEvents(events []*wallet.Event) error
	GetEvents(walletID string) ([]*wallet.Event, error)
	GetEventsAfterVersion(walletID string, version int64) ([]*wallet.Event, error)
	GetAllEvents() ([]*wallet.Event, error)
	CountEvents(walletID string) (int64, error)
}
