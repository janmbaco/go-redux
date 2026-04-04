package application

import "time"

// PlayerDTO represents a player in the application state
type PlayerDTO struct {
	ID       string
	Name     string
	X        float64
	Y        float64
	Health   int
	Score    int
	JoinedAt time.Time
}

// GameState represents the Redux state for the game
type GameState struct {
	Players       map[string]PlayerDTO
	LastEventTime time.Time
}
