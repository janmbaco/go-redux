package domain

import (
	"testing"
)

func TestPlayer_TakeDamage(t *testing.T) {
	tests := []struct {
		name           string
		initialHealth  int
		damage         int
		expectedHealth int
		expectedAlive  bool
	}{
		{"normal damage", 100, 20, 80, true},
		{"damage to zero", 50, 50, 0, false},
		{"overkill damage", 30, 50, 0, false},
		{"zero damage", 100, 0, 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := Player{
				ID:     "test",
				Name:   "Test Player",
				Health: tt.initialHealth,
			}

			player.TakeDamage(tt.damage)

			if player.Health != tt.expectedHealth {
				t.Errorf("Expected health %d, got %d", tt.expectedHealth, player.Health)
			}

			if player.IsAlive() != tt.expectedAlive {
				t.Errorf("Expected alive=%v, got %v", tt.expectedAlive, player.IsAlive())
			}
		})
	}
}

func TestPlayer_Heal(t *testing.T) {
	tests := []struct {
		name           string
		initialHealth  int
		healAmount     int
		expectedHealth int
	}{
		{"normal heal", 50, 30, 80},
		{"heal to max", 80, 30, 100},
		{"overheal capped", 90, 50, 100},
		{"zero heal", 50, 0, 50},
		{"heal from zero", 0, 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			player := Player{
				ID:     "test",
				Name:   "Test Player",
				Health: tt.initialHealth,
			}

			player.Heal(tt.healAmount)

			if player.Health != tt.expectedHealth {
				t.Errorf("Expected health %d, got %d", tt.expectedHealth, player.Health)
			}

			if player.Health > 100 {
				t.Errorf("Health exceeded max: %d", player.Health)
			}
		})
	}
}

func TestPlayer_Move(t *testing.T) {
	player := Player{
		ID:       "test",
		Name:     "Test Player",
		Health:   100,
		Position: Position{X: 0, Y: 0},
	}

	player.Move(5.5, -3.2)

	if player.Position.X != 5.5 {
		t.Errorf("Expected X=5.5, got %f", player.Position.X)
	}

	if player.Position.Y != -3.2 {
		t.Errorf("Expected Y=-3.2, got %f", player.Position.Y)
	}
}
