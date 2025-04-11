package cards

import (
	"github.com/openhs/internal/game"
)

type Jaina struct{}

func (j *Jaina) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Jaina Proudmoore",
		ZhName: "吉安娜·普罗德摩尔",
		ID:     "HERO_08",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
