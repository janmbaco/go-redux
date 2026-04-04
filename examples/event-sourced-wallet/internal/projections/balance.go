package projections

import (
	"sync"

	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"
)

// BalanceProjection is a CQRS read model for wallet balances
type BalanceProjection struct {
	balances map[string]float64 // key: wallet ID
	mu       sync.RWMutex
}

// NewBalanceProjection creates a new balance projection
func NewBalanceProjection() *BalanceProjection {
	return &BalanceProjection{
		balances: make(map[string]float64),
	}
}

// ProcessEvent updates the projection from an event
func (p *BalanceProjection) ProcessEvent(event *wallet.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	switch event.Type {
	case wallet.EventWalletCreated:
		initialBalance := event.Data["initial_balance"].(float64)
		p.balances[event.WalletID] = initialBalance

	case wallet.EventMoneyDeposited:
		amount := event.Data["amount"].(float64)
		p.balances[event.WalletID] += amount

	case wallet.EventMoneyWithdrawn:
		amount := event.Data["amount"].(float64)
		p.balances[event.WalletID] -= amount

	case wallet.EventMoneyTransferred:
		amount := event.Data["amount"].(float64)
		p.balances[event.WalletID] -= amount
		// Note: Transfer creates two events, one for each wallet
	}
}

// GetBalance retrieves current balance for a wallet
func (p *BalanceProjection) GetBalance(walletID string) (float64, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	balance, exists := p.balances[walletID]
	return balance, exists
}

// GetAllBalances retrieves all balances
func (p *BalanceProjection) GetAllBalances() map[string]float64 {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Return copy
	result := make(map[string]float64, len(p.balances))
	for k, v := range p.balances {
		result[k] = v
	}

	return result
}

// Rebuild rebuilds projection from events
func (p *BalanceProjection) Rebuild(events []*wallet.Event) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.balances = make(map[string]float64)

	for _, event := range events {
		p.mu.Unlock() // Unlock temporarily to call ProcessEvent
		p.ProcessEvent(event)
		p.mu.Lock()
	}
}
