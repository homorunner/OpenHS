package engine

import (
	"github.com/openhs/internal/game"
)

func (e *Engine) Attack(attacker *game.Entity, defender *game.Entity, skipValidation bool) error {
	return e.game.Attack(attacker, defender, skipValidation)
}
