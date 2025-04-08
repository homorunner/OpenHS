package engine

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// takes damage
// note: source may be nil
// note: this function will not destroy the entity, that is handled in processGraveyard()
func (e *Engine) DealDamage(source *game.Entity, target *game.Entity, amount int) {
	if target == nil {
		logger.Debug("DealDamage: target is nil, skipping")
		return
	}

	// Deal damage
	target.Health -= amount

	// Check if source has lifesteal
	if source != nil {
		hasLifesteal := game.HasTag(source.Tags, game.TAG_LIFESTEAL)

		// Check if hero with lifesteal weapon
		if !hasLifesteal && source.Card.Type == game.Hero && source.Owner != nil && source.Owner.Weapon != nil {
			hasLifesteal = game.HasTag(source.Owner.Weapon.Tags, game.TAG_LIFESTEAL)
		}

		// Apply lifesteal healing
		if hasLifesteal {
			e.Heal(source, source.Owner.Hero, amount)
		}
	}
}
