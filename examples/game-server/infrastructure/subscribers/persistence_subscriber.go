package subscribers

import (
	"fmt"

	"github.com/janmbaco/go-redux/v2/examples/game-server/application"
)

// PersistenceSubscriber saves state changes to database
type PersistenceSubscriber struct {
	name string
}

// NewPersistenceSubscriber creates a new persistence subscriber
func NewPersistenceSubscriber() *PersistenceSubscriber {
	return &PersistenceSubscriber{
		name: "PersistenceSubscriber",
	}
}

// OnStateChange handles state changes and persists to database
func (s *PersistenceSubscriber) OnStateChange(state application.GameState) {
	fmt.Printf("💾 [PERSISTENCE] Saving game state - %d players\n", len(state.Players))

	// In a real game server, this would save to database
	// For demo, just log the action
}
