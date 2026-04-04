package domain

import (
	"testing"
)

func TestMovementService_MovePlayer(t *testing.T) {
	service := NewMovementService()

	tests := []struct {
		name      string
		fromX     float64
		fromY     float64
		toX       float64
		toY       float64
		expectedX float64
		expectedY float64
	}{
		{"move within bounds", 0, 0, 5, 3, 5, 3},
		{"move from position", 10, 10, 7, 8, 7, 8},
		{"no movement", 5, 5, 5, 5, 5, 5},
		{"negative coords", -5, -3, -3, -2, -3, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newX, newY := service.CalculateNewPosition(tt.fromX, tt.fromY, tt.toX, tt.toY)

			if newX != tt.expectedX {
				t.Errorf("Expected X=%f, got %f", tt.expectedX, newX)
			}

			if newY != tt.expectedY {
				t.Errorf("Expected Y=%f, got %f", tt.expectedY, newY)
			}
		})
	}
}
