package engine

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Heal heals a character by the specified amount
func (e *Engine) Heal(character *game.Entity, amount int) {
	if amount <= 0 {
		logger.Error("Heal: healing amount must be greater than 0")
		return
	}

	// Increase health
	newHealth := character.Health + amount

	// Health cannot exceed MaxHealth
	if newHealth > character.MaxHealth {
		newHealth = character.MaxHealth
	}

	character.Health = newHealth
} 