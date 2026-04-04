package domain

// MovementService handles player movement business rules
type MovementService struct {
	maxSpeed float64
}

// NewMovementService creates a new movement service
func NewMovementService() *MovementService {
	return &MovementService{
		maxSpeed: 10.0,
	}
}

// ValidateMove checks if movement is valid
func (s *MovementService) ValidateMove(fromX, fromY, toX, toY float64) bool {
	dx := toX - fromX
	dy := toY - fromY
	distance := (dx*dx + dy*dy)

	// Check if movement is within max speed
	return distance <= s.maxSpeed*s.maxSpeed
}

// CalculateNewPosition calculates valid new position
func (s *MovementService) CalculateNewPosition(fromX, fromY, toX, toY float64) (float64, float64) {
	if s.ValidateMove(fromX, fromY, toX, toY) {
		return toX, toY
	}

	// Clamp to max speed
	dx := toX - fromX
	dy := toY - fromY
	distance := (dx*dx + dy*dy)

	if distance > s.maxSpeed*s.maxSpeed {
		ratio := s.maxSpeed / distance
		return fromX + dx*ratio, fromY + dy*ratio
	}

	return fromX, fromY
}
