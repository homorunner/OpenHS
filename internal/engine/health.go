package engine

import "github.com/openhs/internal/types"

// SetHealth sets a card's current health and max health to the specified value
func (e *Engine) SetHealth(character *types.Card, health int) {
	character.Health = health
	character.MaxHealth = health
}