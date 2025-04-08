package engine

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Heal heals a character by the specified amount
// note: source may be nil
func (e *Engine) Heal(source *game.Entity, target *game.Entity, amount int) {
	if target == nil {
		logger.Debug("Heal: target is nil, skipping")
		return
	}

	// Increase health
	newHealth := target.Health + amount

	// Health cannot exceed MaxHealth
	if newHealth > target.MaxHealth {
		newHealth = target.MaxHealth
	}

	target.Health = newHealth
}
