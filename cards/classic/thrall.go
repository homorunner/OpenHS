package cards

import (
	"github.com/openhs/internal/game"
)

type Thrall struct{}

func (t *Thrall) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Thrall",
		ZhName: "萨尔",
		ID:     "HERO_02",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
