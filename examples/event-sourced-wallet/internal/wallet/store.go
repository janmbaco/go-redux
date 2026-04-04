package wallet

import (
	"github.com/janmbaco/go-infrastructure/v2/logs"
	"github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/actions"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// WalletState represents the Redux state for wallets
type WalletState struct {
	Wallets map[string]*Wallet `json:"wallets"` // key: wallet ID
	Events  []*Event           `json:"events"`  // Event log
}

// NewWalletState creates initial state
func NewWalletState() WalletState {
	return WalletState{
		Wallets: make(map[string]*Wallet),
		Events:  make([]*Event, 0),
	}
}

// Actions
var (
	CreateWalletAction  = actions.NewAction[WalletCreatedData]("CREATE_WALLET")
	DepositMoneyAction  = actions.NewAction[DepositCommand]("DEPOSIT_MONEY")
	WithdrawMoneyAction = actions.NewAction[WithdrawCommand]("WITHDRAW_MONEY")
	TransferMoneyAction = actions.NewAction[TransferCommand]("TRANSFER_MONEY")
	ReplayEventsAction  = actions.NewAction[[]*Event]("REPLAY_EVENTS")
	ApplySnapshotAction = actions.NewAction[SnapshotCommand]("APPLY_SNAPSHOT")
)

// Command types
type DepositCommand struct {
	WalletID    string  `json:"wallet_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type WithdrawCommand struct {
	WalletID    string  `json:"wallet_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
}

type TransferCommand struct {
	FromWalletID string  `json:"from_wallet_id"`
	ToWalletID   string  `json:"to_wallet_id"`
	Amount       float64 `json:"amount"`
	Description  string  `json:"description"`
}

type SnapshotCommand struct {
	WalletID string        `json:"wallet_id"`
	Snapshot *SnapshotData `json:"snapshot"`
}

// Reducers (pure functions)
func createWalletReducer(state WalletState, data WalletCreatedData) WalletState {
	wallet, err := NewWallet(data.OwnerID, data.Currency, data.InitialBalance)
	if err != nil {
		return state // Invalid command, no state change
	}

	// Create event
	event := &Event{
		ID:        wallet.ID,
		WalletID:  wallet.ID,
		Type:      EventWalletCreated,
		Timestamp: wallet.CreatedAt,
		Version:   1,
		Data: map[string]interface{}{
			"owner_id":        data.OwnerID,
			"currency":        data.Currency,
			"initial_balance": data.InitialBalance,
		},
	}

	state.Wallets[wallet.ID] = wallet
	state.Events = append(state.Events, event)

	return state
}

func depositMoneyReducer(state WalletState, cmd DepositCommand) WalletState {
	wallet, exists := state.Wallets[cmd.WalletID]
	if !exists {
		return state // Wallet not found
	}

	event, err := wallet.Deposit(cmd.Amount, cmd.Description)
	if err != nil {
		return state // Invalid command
	}

	state.Events = append(state.Events, event)
	return state
}

func withdrawMoneyReducer(state WalletState, cmd WithdrawCommand) WalletState {
	wallet, exists := state.Wallets[cmd.WalletID]
	if !exists {
		return state
	}

	event, err := wallet.Withdraw(cmd.Amount, cmd.Description)
	if err != nil {
		return state // Invalid command (insufficient funds)
	}

	state.Events = append(state.Events, event)
	return state
}

func transferMoneyReducer(state WalletState, cmd TransferCommand) WalletState {
	fromWallet, exists := state.Wallets[cmd.FromWalletID]
	if !exists {
		return state
	}

	toWallet, exists := state.Wallets[cmd.ToWalletID]
	if !exists {
		return state
	}

	// Withdraw from source wallet
	withdrawEvent, err := fromWallet.Transfer(cmd.ToWalletID, cmd.Amount, cmd.Description)
	if err != nil {
		return state
	}

	// Deposit to destination wallet
	depositEvent, _ := toWallet.Deposit(cmd.Amount, "Transfer from "+cmd.FromWalletID)

	state.Events = append(state.Events, withdrawEvent, depositEvent)
	return state
}

func replayEventsReducer(state WalletState, events []*Event) WalletState {
	// Clear current state
	newState := NewWalletState()

	// Replay all events
	for _, event := range events {
		switch event.Type {
		case EventWalletCreated:
			data := event.Data
			wallet, _ := NewWallet(
				data["owner_id"].(string),
				data["currency"].(string),
				data["initial_balance"].(float64),
			)
			wallet.ID = event.WalletID
			wallet.Version = event.Version
			wallet.CreatedAt = event.Timestamp
			newState.Wallets[wallet.ID] = wallet

		case EventMoneyDeposited, EventMoneyWithdrawn, EventMoneyTransferred:
			wallet, exists := newState.Wallets[event.WalletID]
			if exists {
				wallet.ApplyEvent(event)
			}
		}

		newState.Events = append(newState.Events, event)
	}

	return newState
}

func applySnapshotReducer(state WalletState, cmd SnapshotCommand) WalletState {
	wallet, exists := state.Wallets[cmd.WalletID]
	if !exists {
		return state
	}

	// Apply snapshot state
	wallet.Balance = cmd.Snapshot.Balance
	wallet.Version = cmd.Snapshot.Version

	// Create snapshot event
	snapshotEvent := wallet.CreateSnapshot(cmd.Snapshot.LastEventID)
	state.Events = append(state.Events, snapshotEvent)

	return state
}

// NewWalletActionHandler creates the Redux action handler
func NewWalletActionHandler() handlers.ActionHandler[map[string]WalletState] {
	builder := handlers.NewActionHandlerBuilder[map[string]WalletState]()

	// Register all reducers
	builder.On(CreateWalletAction, func(state map[string]WalletState, data WalletCreatedData) map[string]WalletState {
		walletState := state["wallets"]
		walletState = createWalletReducer(walletState, data)
		state["wallets"] = walletState
		return state
	})

	builder.On(DepositMoneyAction, func(state map[string]WalletState, cmd DepositCommand) map[string]WalletState {
		walletState := state["wallets"]
		walletState = depositMoneyReducer(walletState, cmd)
		state["wallets"] = walletState
		return state
	})

	builder.On(WithdrawMoneyAction, func(state map[string]WalletState, cmd WithdrawCommand) map[string]WalletState {
		walletState := state["wallets"]
		walletState = withdrawMoneyReducer(walletState, cmd)
		state["wallets"] = walletState
		return state
	})

	builder.On(TransferMoneyAction, func(state map[string]WalletState, cmd TransferCommand) map[string]WalletState {
		walletState := state["wallets"]
		walletState = transferMoneyReducer(walletState, cmd)
		state["wallets"] = walletState
		return state
	})

	builder.On(ReplayEventsAction, func(state map[string]WalletState, events []*Event) map[string]WalletState {
		walletState := state["wallets"]
		walletState = replayEventsReducer(walletState, events)
		state["wallets"] = walletState
		return state
	})

	builder.On(ApplySnapshotAction, func(state map[string]WalletState, cmd SnapshotCommand) map[string]WalletState {
		walletState := state["wallets"]
		walletState = applySnapshotReducer(walletState, cmd)
		state["wallets"] = walletState
		return state
	})

	return builder.Build()
}

// CreateWalletStore creates a new Redux store for wallets
func CreateWalletStore(logger logs.Logger) redux.Store[map[string]WalletState] {
	initialState := map[string]WalletState{
		"wallets": NewWalletState(),
	}

	store := redux.NewStore(initialState, logger)
	store.AddModule(NewWalletActionHandler())

	return store
}
