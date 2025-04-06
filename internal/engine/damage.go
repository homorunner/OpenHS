package engine

import "github.com/openhs/internal/game"

// takes damage
// note: this function will not destroy the entity, that is handled in processGraveyard()
func (e *Engine) TakeDamage(character *game.Entity, amount int) {
	character.Health -= amount
}

