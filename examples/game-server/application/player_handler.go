package application

import (
	"time"

	"github.com/janmbaco/go-redux/v2/examples/game-server/domain"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// PlayerHandler handles player-related actions
type PlayerHandler struct {
	movementService *domain.MovementService
}

// NewPlayerHandler creates a new player handler
func NewPlayerHandler(movementService *domain.MovementService) handlers.ActionHandler[GameState] {
	handler := &PlayerHandler{
		movementService: movementService,
	}

	return handlers.NewActionHandlerBuilder[GameState]().
		On(JoinGameAction, handler.handleJoinGame).
		On(LeaveGameAction, handler.handleLeaveGame).
		On(MovePlayerAction, handler.handleMovePlayer).
		Build()
}

func (h *PlayerHandler) handleJoinGame(state GameState, payload JoinGamePayload) GameState {
	// Create new player using domain
	player := domain.NewPlayer(payload.PlayerID, payload.Name)

	// Convert to DTO and add to state
	state.Players[player.ID] = ToDTO(player)
	state.LastEventTime = time.Now()

	return state
}

func (h *PlayerHandler) handleLeaveGame(state GameState, playerID string) GameState {
	// Remove player from state
	delete(state.Players, playerID)
	state.LastEventTime = time.Now()

	return state
}

func (h *PlayerHandler) handleMovePlayer(state GameState, payload MovePlayerPayload) GameState {
	player, exists := state.Players[payload.PlayerID]
	if !exists {
		return state
	}

	// Use domain service to validate and calculate movement
	newX, newY := h.movementService.CalculateNewPosition(
		player.X, player.Y,
		payload.X, payload.Y,
	)

	// Update position
	player.X = newX
	player.Y = newY
	state.Players[payload.PlayerID] = player
	state.LastEventTime = time.Now()

	return state
}
