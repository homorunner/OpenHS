package engine

import "github.com/openhs/internal/game"

// SetHealth sets a card's current health and max health to the specified value
func (e *Engine) SetHealth(character *game.Card, health int) {
	character.Health = health
	character.MaxHealth = health
}