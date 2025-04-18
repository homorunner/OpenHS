package cards

import (
	"github.com/openhs/internal/game"
)

type Frostbolt struct{}

func (f *Frostbolt) Register(cm *game.CardManager) {
	card := game.Card{
		Name:        "Frostbolt",
		ZhName:      "寒冰箭",
		ID:          "CORE_CS2_024",
		Description: "对一个角色造成3点伤害，并使其冻结。",
		Cost:        2,
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

func (f *Frostbolt) Cast(g *game.Game, source *game.Entity, target *game.Entity) {
	// Deal 3 damage to the target
	g.DealDamage(source, target, 3)

	// Freeze the target
	g.Freeze(target)
}
