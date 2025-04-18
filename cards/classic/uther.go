package cards

import (
	"github.com/openhs/internal/game"
)

type Uther struct{}

func (u *Uther) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Uther Lightbringer",
		ZhName: "乌瑟尔·光明使者",
		ID:     "HERO_04",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
