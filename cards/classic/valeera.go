package cards

import (
	"github.com/openhs/internal/game"
)

type Valeera struct{}

func (v *Valeera) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Valeera Sanguinar",
		ZhName: "瓦莉拉·萨古纳尔",
		ID:     "HERO_03",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
