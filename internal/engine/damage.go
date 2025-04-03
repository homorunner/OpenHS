package engine

import "github.com/openhs/internal/game"

func (e *Engine) TakeDamage(character *game.Card, amount int) {
	character.Health -= amount
}

