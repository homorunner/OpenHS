package engine

import "github.com/openhs/internal/types"

func (e *Engine) TakeDamage(character *types.Card, amount int) {
	character.Health -= amount
}

