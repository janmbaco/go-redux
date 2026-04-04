package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/janmbaco/go-infrastructure/v2/logs"
	redux "github.com/janmbaco/go-redux/v2"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/eventstore"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/projections"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/snapshots"
	"github.com/janmbaco/go-redux/v2/examples/event-sourced-wallet/internal/wallet"
)

type WalletHandler struct {
	store       redux.Store[map[string]wallet.WalletState]
	eventStore  eventstore.EventStore
	projection  *projections.BalanceProjection
	snapshotter *snapshots.Snapshotter
	logger      logs.Logger
}

func NewWalletHandler(
	store redux.Store[map[string]wallet.WalletState],
	eventStore eventstore.EventStore,
	projection *projections.BalanceProjection,
	snapshotter *snapshots.Snapshotter,
	logger logs.Logger,
) *WalletHandler {
	return &WalletHandler{
		store:       store,
		eventStore:  eventStore,
		projection:  projection,
		snapshotter: snapshotter,
		logger:      logger,
	}
}

// CreateWallet handles POST /wallets
func (h *WalletHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		OwnerID        string  `json:"owner_id"`
		Currency       string  `json:"currency"`
		InitialBalance float64 `json:"initial_balance"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	// Dispatch CREATE_WALLET action
	data := wallet.WalletCreatedData{
		OwnerID:        req.OwnerID,
		Currency:       req.Currency,
		InitialBalance: req.InitialBalance,
	}

	if err := h.store.Dispatch(wallet.CreateWalletAction.With(data)); err != nil {
		h.logger.Error("[WalletHandler] Failed to create wallet: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create wallet"})
		return
	}

	// Get created wallet from state
	state := h.store.GetState()
	walletState := state["wallets"]

	// Find the newly created wallet
	var createdWallet *wallet.Wallet
	for _, w := range walletState.Wallets {
		if w.OwnerID == req.OwnerID && w.Currency == req.Currency {
			createdWallet = w
			break
		}
	}

	if createdWallet == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "wallet creation failed"})
		return
	}

	// Save event to event store
	if len(walletState.Events) > 0 {
		lastEvent := walletState.Events[len(walletState.Events)-1]
		h.eventStore.SaveEvent(lastEvent)
		h.projection.ProcessEvent(lastEvent)
	}

	h.logger.Info("[WalletHandler] Wallet created: " + createdWallet.ID)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWallet)
}

// Deposit handles POST /wallets/:id/deposit
func (h *WalletHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("id")

	var req struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	cmd := wallet.DepositCommand{
		WalletID:    walletID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	if err := h.store.Dispatch(wallet.DepositMoneyAction.With(cmd)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to deposit"})
		return
	}

	// Save event
	state := h.store.GetState()
	walletState := state["wallets"]
	if len(walletState.Events) > 0 {
		lastEvent := walletState.Events[len(walletState.Events)-1]
		h.eventStore.SaveEvent(lastEvent)
		h.projection.ProcessEvent(lastEvent)

		// Check snapshot
		eventCount, _ := h.eventStore.CountEvents(walletID)
		if w, exists := walletState.Wallets[walletID]; exists {
			h.snapshotter.CheckAndCreateSnapshot(w, eventCount, lastEvent.ID)
		}
	}

	h.logger.Info("[WalletHandler] Deposited to wallet: " + walletID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "deposit successful"})
}

// Withdraw handles POST /wallets/:id/withdraw
func (h *WalletHandler) Withdraw(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("id")

	var req struct {
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	cmd := wallet.WithdrawCommand{
		WalletID:    walletID,
		Amount:      req.Amount,
		Description: req.Description,
	}

	if err := h.store.Dispatch(wallet.WithdrawMoneyAction.With(cmd)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to withdraw"})
		return
	}

	// Save event
	state := h.store.GetState()
	walletState := state["wallets"]
	if len(walletState.Events) > 0 {
		lastEvent := walletState.Events[len(walletState.Events)-1]
		h.eventStore.SaveEvent(lastEvent)
		h.projection.ProcessEvent(lastEvent)

		// Check snapshot
		eventCount, _ := h.eventStore.CountEvents(walletID)
		if w, exists := walletState.Wallets[walletID]; exists {
			h.snapshotter.CheckAndCreateSnapshot(w, eventCount, lastEvent.ID)
		}
	}

	h.logger.Info("[WalletHandler] Withdrawn from wallet: " + walletID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "withdrawal successful"})
}

// Transfer handles POST /wallets/:id/transfer
func (h *WalletHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	fromWalletID := r.PathValue("id")

	var req struct {
		ToWalletID  string  `json:"to_wallet_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	cmd := wallet.TransferCommand{
		FromWalletID: fromWalletID,
		ToWalletID:   req.ToWalletID,
		Amount:       req.Amount,
		Description:  req.Description,
	}

	if err := h.store.Dispatch(wallet.TransferMoneyAction.With(cmd)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to transfer"})
		return
	}

	// Save events (transfer creates 2 events)
	state := h.store.GetState()
	walletState := state["wallets"]

	if len(walletState.Events) >= 2 {
		// Get last 2 events
		event1 := walletState.Events[len(walletState.Events)-2]
		event2 := walletState.Events[len(walletState.Events)-1]

		h.eventStore.SaveEvent(event1)
		h.eventStore.SaveEvent(event2)
		h.projection.ProcessEvent(event1)
		h.projection.ProcessEvent(event2)
	}

	h.logger.Info("[WalletHandler] Transfer from " + fromWalletID + " to " + req.ToWalletID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "transfer successful"})
}

// GetWallet handles GET /wallets/:id
func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("id")

	state := h.store.GetState()
	walletState := state["wallets"]

	wallet, exists := walletState.Wallets[walletID]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "wallet not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

// GetBalance handles GET /wallets/:id/balance (from projection)
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("id")

	balance, exists := h.projection.GetBalance(walletID)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "wallet not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"wallet_id": walletID,
		"balance":   balance,
	})
}

// GetEvents handles GET /wallets/:id/events
func (h *WalletHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("id")

	events, err := h.eventStore.GetEvents(walletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get events"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"wallet_id":   walletID,
		"event_count": len(events),
		"events":      events,
	})
}

// Replay handles POST /wallets/:id/replay
func (h *WalletHandler) Replay(w http.ResponseWriter, r *http.Request) {
	walletID := r.PathValue("id")

	events, err := h.eventStore.GetEvents(walletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to get events"})
		return
	}

	h.logger.Info("[WalletHandler] Replaying " + string(rune(len(events))) + " events for wallet " + walletID)

	// Dispatch replay action
	if err := h.store.Dispatch(wallet.ReplayEventsAction.With(events)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to replay"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":      "replay completed",
		"events_count": len(events),
	})
}

// HealthCheck handles GET /health
func (h *WalletHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"service": "event-sourced-wallet",
		"status":  "healthy",
	})
}
