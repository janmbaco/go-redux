package wallet

import "time"

// EventType represents the type of wallet event
type EventType string

const (
	EventWalletCreated    EventType = "WALLET_CREATED"
	EventMoneyDeposited   EventType = "MONEY_DEPOSITED"
	EventMoneyWithdrawn   EventType = "MONEY_WITHDRAWN"
	EventMoneyTransferred EventType = "MONEY_TRANSFERRED"
	EventSnapshotTaken    EventType = "SNAPSHOT_TAKEN"
)

// Event represents a domain event in the wallet
type Event struct {
	ID            string                 `json:"id"`
	WalletID      string                 `json:"wallet_id"`
	Type          EventType              `json:"type"`
	Data          map[string]interface{} `json:"data"`
	Timestamp     time.Time              `json:"timestamp"`
	Version       int64                  `json:"version"` // Aggregate version for optimistic locking
	CorrelationID string                 `json:"correlation_id"`
}

// WalletCreatedData contains data for wallet creation event
type WalletCreatedData struct {
	OwnerID        string  `json:"owner_id"`
	Currency       string  `json:"currency"`
	InitialBalance float64 `json:"initial_balance"`
}

// MoneyDepositedData contains data for deposit event
type MoneyDepositedData struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// MoneyWithdrawnData contains data for withdrawal event
type MoneyWithdrawnData struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// MoneyTransferredData contains data for transfer event
type MoneyTransferredData struct {
	ToWalletID  string  `json:"to_wallet_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

// SnapshotData contains wallet snapshot state
type SnapshotData struct {
	Balance     float64 `json:"balance"`
	Version     int64   `json:"version"`
	LastEventID string  `json:"last_event_id"`
}
