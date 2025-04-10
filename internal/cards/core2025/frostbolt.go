package cards

import (
	"github.com/openhs/internal/game"
)

type Frostbolt struct{}

func (f *Frostbolt) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Frostbolt",
		ZhName: "寒冰箭",
		ID:     "CORE_CS2_024",
		Cost:   2,
		Type:   game.Spell,
	}

	cm.RegisterCard(card)
}
