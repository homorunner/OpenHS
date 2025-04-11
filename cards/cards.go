package cards

import (
	core "github.com/openhs/cards/core2025"
	edr "github.com/openhs/cards/edr"
	"github.com/openhs/internal/game"
)

type CardDef interface {
	Register(cm *game.CardManager)
}

func RegisterAllCards(cm *game.CardManager) {
	for _, card := range core.AllCards {
		card.(CardDef).Register(cm)
	}
	for _, card := range edr.AllCards {
		card.(CardDef).Register(cm)
	}
}
