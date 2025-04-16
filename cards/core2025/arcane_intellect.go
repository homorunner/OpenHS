package cards

import (
	"github.com/openhs/internal/game"
)

type ArcaneIntellect struct{}

func (a *ArcaneIntellect) Register(cm *game.CardManager) {
	card := game.Card{
		Name:        "Arcane Intellect",
		ZhName:      "奥术智慧",
		ID:          "CORE_CS2_023",
		Description: "抽两张牌。",
		Cost:        3,
		Type:        game.Spell,
		Powers: []game.Power{
			{
				Type:   game.PowerTypeSpell,
				Action: a.Cast,
			},
		},
	}

	cm.RegisterCard(card)
}

func (a *ArcaneIntellect) Cast(g *game.Game, source *game.Entity, target *game.Entity) {
	drawCount := 2
	for i := 0; i < drawCount; i++ {
		if source.Owner != nil {
			g.DrawCard(source.Owner)
		}
	}
}
