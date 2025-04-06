package engine

import "github.com/openhs/internal/game"

func (e *Engine) removeEntityFromBoard(player *game.Player, entity *game.Entity) {
	for i, minion := range player.Field {
		if minion == entity {
			player.Field = append(player.Field[:i], player.Field[i+1:]...)
			break
		}
	}
}

