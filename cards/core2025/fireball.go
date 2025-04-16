package cards

import (
	"github.com/openhs/internal/game"
)

type Fireball struct{}

func (f *Fireball) Register(cm *game.CardManager) {
	card := game.Card{
		Name:        "Fireball",
		ZhName:      "火球术",
		ID:          "CORE_CS2_029",
		Description: "造成6点伤害。",
		Cost:        4,
		Type:        game.Spell,
		Powers: []game.Power{
			{
				Type:   game.PowerTypeSpell,
				Action: f.Cast,
			},
		},
	}

	cm.RegisterCard(card)
}

func (f *Fireball) Cast(g *game.Game, source *game.Entity, target *game.Entity) {
	g.DealDamage(source, target, 6)
}
