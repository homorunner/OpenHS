package engine

import (
	"github.com/openhs/internal/game"
)

func (e *Engine) removeEntityFromBoard(player *game.Player, entity *game.Entity) {
	for i, minion := range player.Field {
		if minion == entity {
			player.Field = append(player.Field[:i], player.Field[i+1:]...)
			// The entity is being removed from the board,
			// but it will be placed in a new zone afterwards,
			// so we set it to NONE temporarily
			entity.CurrentZone = game.ZONE_NONE
			break
		}
	}
}
