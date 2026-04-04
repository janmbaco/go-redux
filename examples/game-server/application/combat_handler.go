package application

import (
	"time"

	"github.com/janmbaco/go-redux/v2/examples/game-server/domain"
	"github.com/janmbaco/go-redux/v2/handlers"
)

// CombatHandler handles combat-related actions
type CombatHandler struct {
	combatService *domain.CombatService
}

// NewCombatHandler creates a new combat handler
func NewCombatHandler(combatService *domain.CombatService) handlers.ActionHandler[GameState] {
	handler := &CombatHandler{
		combatService: combatService,
	}

	return handlers.NewActionHandlerBuilder[GameState]().
		On(AttackPlayerAction, handler.handleAttack).
		On(HealPlayerAction, handler.handleHeal).
		Build()
}

func (h *CombatHandler) handleAttack(state GameState, payload AttackPayload) GameState {
	attacker, attackerExists := state.Players[payload.AttackerID]
	victim, victimExists := state.Players[payload.VictimID]

	if !attackerExists || !victimExists {
		return state
	}

	// Calculate damage using domain service
	damage := h.combatService.CalculateDamage(attacker.Score)

	// Apply damage
	victim.Health -= damage
	if victim.Health < 0 {
		victim.Health = 0
	}

	// If victim died, give attacker score reward
	if victim.Health == 0 {
		reward := h.combatService.CalculateScoreReward(victim.Score)
		attacker.Score += reward
		state.Players[payload.AttackerID] = attacker
	}

	state.Players[payload.VictimID] = victim
	state.LastEventTime = time.Now()

	return state
}

func (h *CombatHandler) handleHeal(state GameState, payload HealPayload) GameState {
	player, exists := state.Players[payload.PlayerID]
	if !exists {
		return state
	}

	// Heal player
	player.Health += payload.Amount
	if player.Health > 100 {
		player.Health = 100
	}

	state.Players[payload.PlayerID] = player
	state.LastEventTime = time.Now()

	return state
}
