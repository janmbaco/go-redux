package subscribers

import (
	"fmt"

	"github.com/janmbaco/go-redux/v2/examples/game-server/application"
)

// AnalyticsSubscriber tracks game analytics
type AnalyticsSubscriber struct {
	name string
}

// NewAnalyticsSubscriber creates a new analytics subscriber
func NewAnalyticsSubscriber() *AnalyticsSubscriber {
	return &AnalyticsSubscriber{
		name: "AnalyticsSubscriber",
	}
}

// OnStateChange handles state changes and tracks analytics
func (s *AnalyticsSubscriber) OnStateChange(state application.GameState) {
	// Calculate analytics
	totalHealth := 0
	totalScore := 0
	alivePlayers := 0

	for _, player := range state.Players {
		totalHealth += player.Health
		totalScore += player.Score
		if player.Health > 0 {
			alivePlayers++
		}
	}

	if len(state.Players) > 0 {
		avgHealth := totalHealth / len(state.Players)
		avgScore := totalScore / len(state.Players)

		fmt.Printf("📊 [ANALYTICS] Avg Health: %d, Avg Score: %d, Alive: %d/%d\n",
			avgHealth, avgScore, alivePlayers, len(state.Players))
	}
}
