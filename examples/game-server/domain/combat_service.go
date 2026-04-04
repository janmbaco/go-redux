package domain

// CombatService handles combat business rules
type CombatService struct {
	baseDamage   int
	criticalRate float64
}

// NewCombatService creates a new combat service
func NewCombatService() *CombatService {
	return &CombatService{
		baseDamage:   10,
		criticalRate: 0.2, // 20% critical hit chance
	}
}

// CalculateDamage calculates damage dealt
func (s *CombatService) CalculateDamage(attackerScore int) int {
	damage := s.baseDamage

	// Bonus damage based on score
	if attackerScore > 100 {
		damage += 5
	}
	if attackerScore > 500 {
		damage += 10
	}

	return damage
}

// CalculateScoreReward calculates score for killing a player
func (s *CombatService) CalculateScoreReward(victimScore int) int {
	// Base reward: 10 points
	reward := 10

	// Bonus for killing high-score players
	if victimScore > 100 {
		reward += 20
	}
	if victimScore > 500 {
		reward += 50
	}

	return reward
}

// ShouldRespawn checks if player should respawn
func (s *CombatService) ShouldRespawn(timeSinceDeath float64) bool {
	// Respawn after 5 seconds
	return timeSinceDeath >= 5.0
}
