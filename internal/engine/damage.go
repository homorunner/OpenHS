package engine

import (
	"github.com/openhs/internal/game"
)

// takes damage
// note: source may be nil
// note: this function will not destroy the entity, that is handled in processGraveyard()
func (e *Engine) DealDamage(source *game.Entity, target *game.Entity, amount int) {
	e.game.DealDamage(source, target, amount)
}
