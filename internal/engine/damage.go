package engine

import "github.com/openhs/internal/game"

func (e *Engine) TakeDamage(character *game.Entity, amount int) {
	character.Health -= amount
}

