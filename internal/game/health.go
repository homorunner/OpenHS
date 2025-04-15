package game

import (
	"github.com/openhs/internal/logger"
)

// DealDamage deals damage to a target from a source
// note: source may be nil
// note: this function will not destroy the entity, that is handled elsewhere
func (g *Game) DealDamage(source *Entity, target *Entity, amount int) {
	if target == nil {
		logger.Debug("DealDamage: target is nil, skipping")
		return
	}

	if amount <= 0 {
		logger.Debug("DealDamage: amount is <= 0, skipping")
		return
	}

	// Create context for damage trigger
	damageCtx := TriggerContext{
		Game:         g,
		SourceEntity: source,
		TargetEntity: target,
		Value:        amount,
		Phase:        g.Phase,
	}

	// Deal damage
	target.Health -= amount

	// Trigger damage taken event
	g.TriggerManager.ActivateTrigger(TriggerDamageTaken, damageCtx)

	// Also trigger hero damage taken event if the target is a hero
	if target.Card.Type == Hero {
		g.TriggerManager.ActivateTrigger(TriggerHeroDamageTaken, damageCtx)
	}

	// Check if source has lifesteal
	if source != nil {
		hasLifesteal := HasTag(source.Tags, TAG_LIFESTEAL)

		// Check if hero with lifesteal weapon
		if !hasLifesteal && source.Card.Type == Hero && source.Owner != nil && source.Owner.Weapon != nil {
			hasLifesteal = HasTag(source.Owner.Weapon.Tags, TAG_LIFESTEAL)
		}

		// Apply lifesteal healing
		if hasLifesteal {
			g.Heal(source, source.Owner.Hero, amount)
		}
	}
}

// Heal heals a character by the specified amount
// note: source may be nil
func (g *Game) Heal(source *Entity, target *Entity, amount int) {
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
	healCtx := TriggerContext{
		Game:         g,
		SourceEntity: source,
		TargetEntity: target,
		Value:        healedAmount,
		Phase:        g.Phase,
	}

	// Trigger heal received event
	g.TriggerManager.ActivateTrigger(TriggerHealReceived, healCtx)
}

// SetHealth sets a card's current health and max health to the specified value
func (g *Game) SetHealth(character *Entity, health int) {
	character.Health = health
	character.MaxHealth = health
} 