package cards

import (
	"github.com/openhs/internal/game"
)

type Garrosh struct{}

func (g *Garrosh) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Garrosh Hellscream",
		ZhName: "加尔鲁什·地狱咆哮",
		ID:     "HERO_01",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
