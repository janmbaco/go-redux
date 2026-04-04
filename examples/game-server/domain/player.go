package domain

import "time"

// Player represents a player entity
type Player struct {
	ID       string
	Name     string
	Position Position
	Health   int
	Score    int
	JoinedAt time.Time
}

// Position represents a 2D position
type Position struct {
	X float64
	Y float64
}

// NewPlayer creates a new player
func NewPlayer(id, name string) *Player {
	return &Player{
		ID:       id,
		Name:     name,
		Position: Position{X: 0, Y: 0},
		Health:   100,
		Score:    0,
		JoinedAt: time.Now(),
	}
}

// IsAlive checks if player is alive
func (p *Player) IsAlive() bool {
	return p.Health > 0
}

// TakeDamage applies damage to player
func (p *Player) TakeDamage(amount int) {
	p.Health -= amount
	if p.Health < 0 {
		p.Health = 0
	}
}

// Heal restores health
func (p *Player) Heal(amount int) {
	p.Health += amount
	if p.Health > 100 {
		p.Health = 100
	}
}

// Move changes player position
func (p *Player) Move(x, y float64) {
	p.Position.X = x
	p.Position.Y = y
}

// AddScore increases player score
func (p *Player) AddScore(points int) {
	p.Score += points
}
