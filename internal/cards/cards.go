package cards

import (
	core "github.com/openhs/internal/cards/core2025"
	edr "github.com/openhs/internal/cards/edr"
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
