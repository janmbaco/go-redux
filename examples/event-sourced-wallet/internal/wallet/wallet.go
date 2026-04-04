package wallet

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("invalid amount: must be positive")
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInvalidCurrency   = errors.New("invalid currency")
)

// Wallet represents a wallet aggregate
type Wallet struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"owner_id"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	Version   int64     `json:"version"` // For event sourcing
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewWallet creates a new wallet
func NewWallet(ownerID, currency string, initialBalance float64) (*Wallet, error) {
	if initialBalance < 0 {
		return nil, ErrInvalidAmount
	}

	if currency == "" {
		return nil, ErrInvalidCurrency
	}

	now := time.Now()
	return &Wallet{
		ID:        uuid.New().String(),
		OwnerID:   ownerID,
		Currency:  currency,
		Balance:   initialBalance,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Deposit adds money to the wallet (command)
func (w *Wallet) Deposit(amount float64, description string) (*Event, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	w.Balance += amount
	w.Version++
	w.UpdatedAt = time.Now()

	event := &Event{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      EventMoneyDeposited,
		Timestamp: time.Now(),
		Version:   w.Version,
		Data: map[string]interface{}{
			"amount":      amount,
			"description": description,
		},
	}

	return event, nil
}

// Withdraw removes money from the wallet (command)
func (w *Wallet) Withdraw(amount float64, description string) (*Event, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if w.Balance < amount {
		return nil, ErrInsufficientFunds
	}

	w.Balance -= amount
	w.Version++
	w.UpdatedAt = time.Now()

	event := &Event{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      EventMoneyWithdrawn,
		Timestamp: time.Now(),
		Version:   w.Version,
		Data: map[string]interface{}{
			"amount":      amount,
			"description": description,
		},
	}

	return event, nil
}

// Transfer sends money to another wallet (command)
func (w *Wallet) Transfer(toWalletID string, amount float64, description string) (*Event, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	if w.Balance < amount {
		return nil, ErrInsufficientFunds
	}

	w.Balance -= amount
	w.Version++
	w.UpdatedAt = time.Now()

	event := &Event{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      EventMoneyTransferred,
		Timestamp: time.Now(),
		Version:   w.Version,
		Data: map[string]interface{}{
			"to_wallet_id": toWalletID,
			"amount":       amount,
			"description":  description,
		},
	}

	return event, nil
}

// ApplyEvent applies an event to rebuild wallet state (for event sourcing replay)
func (w *Wallet) ApplyEvent(event *Event) error {
	switch event.Type {
	case EventWalletCreated:
		data := event.Data
		w.ID = event.WalletID
		w.OwnerID = data["owner_id"].(string)
		w.Currency = data["currency"].(string)
		w.Balance = data["initial_balance"].(float64)
		w.CreatedAt = event.Timestamp

	case EventMoneyDeposited:
		amount := event.Data["amount"].(float64)
		w.Balance += amount

	case EventMoneyWithdrawn:
		amount := event.Data["amount"].(float64)
		w.Balance -= amount

	case EventMoneyTransferred:
		amount := event.Data["amount"].(float64)
		w.Balance -= amount

	default:
		return errors.New("unknown event type")
	}

	w.Version = event.Version
	w.UpdatedAt = event.Timestamp

	return nil
}

// CreateSnapshot generates a snapshot of current wallet state
func (w *Wallet) CreateSnapshot(lastEventID string) *Event {
	w.Version++

	return &Event{
		ID:        uuid.New().String(),
		WalletID:  w.ID,
		Type:      EventSnapshotTaken,
		Timestamp: time.Now(),
		Version:   w.Version,
		Data: map[string]interface{}{
			"balance":       w.Balance,
			"version":       w.Version,
			"last_event_id": lastEventID,
		},
	}
}
