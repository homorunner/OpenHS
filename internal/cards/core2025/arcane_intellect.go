package cards

import "github.com/openhs/internal/game"

type ArcaneIntellect struct{}

func (a *ArcaneIntellect) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Arcane Intellect",
		ZhName: "奥术智慧",
		ID:     "CORE_CS2_023",
		Cost:   3,
		Type:   game.Spell,
	}

	cm.RegisterCard(card)
}
