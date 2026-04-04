package domain

import (
	"testing"
)

func TestCombatService_CalculateDamage(t *testing.T) {
	service := NewCombatService()

	tests := []struct {
		name          string
		attackerScore int
		expectedMin   int
	}{
		{"base damage", 0, 10},
		{"with bonus", 150, 15},
		{"high score", 600, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			damage := service.CalculateDamage(tt.attackerScore)

			if damage < tt.expectedMin {
				t.Errorf("Expected damage >= %d, got %d", tt.expectedMin, damage)
			}
		})
	}
}

func TestCombatService_CalculateScoreReward(t *testing.T) {
	service := NewCombatService()

	reward := service.CalculateScoreReward(50)
	if reward < 10 {
		t.Errorf("Expected reward >= 10, got %d", reward)
	}

	highReward := service.CalculateScoreReward(150)
	if highReward < 30 {
		t.Errorf("Expected high reward >= 30, got %d", highReward)
	}
}
