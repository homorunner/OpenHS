package cards

import (
	"github.com/openhs/internal/game"
)

type Fireball struct{}

func (f *Fireball) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Fireball",
		ZhName: "火球术",
		ID:     "CORE_CS2_029",
		Cost:   4,
		Type:   game.Spell,
	}

	cm.RegisterCard(card)
}
