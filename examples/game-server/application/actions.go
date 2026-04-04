package application

import "github.com/janmbaco/go-redux/v2/actions"

// Player actions
var JoinGameAction = actions.NewAction[JoinGamePayload]("game/joinGame")
var LeaveGameAction = actions.NewAction[string]("game/leaveGame")
var MovePlayerAction = actions.NewAction[MovePlayerPayload]("game/movePlayer")
var AttackPlayerAction = actions.NewAction[AttackPayload]("game/attackPlayer")
var HealPlayerAction = actions.NewAction[HealPayload]("game/healPlayer")

// Payloads
type JoinGamePayload struct {
	PlayerID string
	Name     string
}

type MovePlayerPayload struct {
	PlayerID string
	X        float64
	Y        float64
}

type AttackPayload struct {
	AttackerID string
	VictimID   string
}

type HealPayload struct {
	PlayerID string
	Amount   int
}
