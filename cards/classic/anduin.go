package cards

import (
	"github.com/openhs/internal/game"
)

type Anduin struct{}

func (a *Anduin) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Anduin Wrynn",
		ZhName: "安度因·乌瑞恩",
		ID:     "HERO_09",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
