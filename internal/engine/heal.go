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

	if amount <= 0 {
		logger.Debug("Heal: amount is <= 0, skipping")
		return
	}

	// Increase health
	oldHealth := target.Health
	newHealth := oldHealth + amount

	// Health cannot exceed MaxHealth
	if newHealth > target.MaxHealth {
		newHealth = target.MaxHealth
	}

	// Calculate actual amount healed (could be less than amount if hitting max health)
	healedAmount := newHealth - oldHealth

	// If no actual healing occurred, skip
	if healedAmount <= 0 {
		return
	}

	// Apply the heal
	target.Health = newHealth

	// Create context for heal trigger
	healCtx := game.TriggerContext{
		Game:         e.game,
		SourceEntity: source,
		TargetEntity: target,
		Value:        healedAmount,
		Phase:        e.game.Phase,
	}

	// Trigger heal received event
	e.game.TriggerManager.ActivateTrigger(game.TriggerHealReceived, healCtx)
}
