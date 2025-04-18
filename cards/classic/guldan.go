package cards

import (
	"github.com/openhs/internal/game"
)

type Guldan struct{}

func (g *Guldan) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Gul'dan",
		ZhName: "古尔丹",
		ID:     "HERO_07",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
