package subscribers

import (
	"fmt"

	"github.com/janmbaco/go-redux/v2/examples/game-server/application"
)

// BroadcastSubscriber broadcasts state changes to all connected clients
type BroadcastSubscriber struct {
	name string
}

// NewBroadcastSubscriber creates a new broadcast subscriber
func NewBroadcastSubscriber() *BroadcastSubscriber {
	return &BroadcastSubscriber{
		name: "BroadcastSubscriber",
	}
}

// OnStateChange handles state changes and broadcasts to clients
func (s *BroadcastSubscriber) OnStateChange(state application.GameState) {
	fmt.Printf("📡 [BROADCAST] State update - Players: %d, Last event: %s\n",
		len(state.Players),
		state.LastEventTime.Format("15:04:05"),
	)

	// In a real game server, this would broadcast to WebSocket clients
	for id, player := range state.Players {
		fmt.Printf("   Player %s: %s [HP:%d Score:%d Pos:(%.1f,%.1f)]\n",
			id, player.Name, player.Health, player.Score, player.X, player.Y)
	}
}
