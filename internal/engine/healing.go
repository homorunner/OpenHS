package engine

import (
	"github.com/openhs/internal/logger"
	"github.com/openhs/internal/types"
)

// HealCard heals a card by the specified amount
func (e *Engine) HealCard(character *types.Card, amount int) {
	if amount <= 0 {
		logger.Error("HealCard: healing amount must be greater than 0")
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