package application

import (
	"github.com/janmbaco/go-redux/v2/examples/game-server/domain"
)

// ToDTO converts a domain Player to DTO
func ToDTO(player *domain.Player) PlayerDTO {
	return PlayerDTO{
		ID:       player.ID,
		Name:     player.Name,
		X:        player.Position.X,
		Y:        player.Position.Y,
		Health:   player.Health,
		Score:    player.Score,
		JoinedAt: player.JoinedAt,
	}
}

// ToEntity converts DTO to domain Player
func ToEntity(dto PlayerDTO) *domain.Player {
	player := &domain.Player{
		ID:       dto.ID,
		Name:     dto.Name,
		Position: domain.Position{X: dto.X, Y: dto.Y},
		Health:   dto.Health,
		Score:    dto.Score,
		JoinedAt: dto.JoinedAt,
	}
	return player
}
