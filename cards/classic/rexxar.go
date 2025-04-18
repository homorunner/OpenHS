package cards

import (
	"github.com/openhs/internal/game"
)

type Rexxar struct{}

func (r *Rexxar) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Rexxar",
		ZhName: "雷克萨",
		ID:     "HERO_05",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
