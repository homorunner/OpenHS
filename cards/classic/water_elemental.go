package cards

import (
	"github.com/openhs/internal/game"
)

type WaterElemental struct{}

func (w *WaterElemental) Register(cm *game.CardManager) {
	card := game.Card{
		Name:        "Water Elemental",
		ZhName:      "水元素",
		ID:          "CS2_033",
		Description: "冻结任何受到本随从伤害的角色。",
		Cost:        4,
		Attack:      3,
		Health:      6,
		Type:        game.Minion,
		Load:        w.Load,
		Unload:      w.Unload,
	}

	cm.RegisterCard(card)
}

func (w *WaterElemental) Load(g *game.Game, self *game.Entity) {
	g.TriggerManager.RegisterTrigger(game.TriggerDamageTaken, self, w.OnDamage, false)
}

func (w *WaterElemental) Unload(g *game.Game, self *game.Entity) {
	g.TriggerManager.UnregisterAllForEntity(self)
}

func (w *WaterElemental) OnDamage(ctx *game.TriggerContext, self *game.Entity) {
	if ctx.SourceEntity != self {
		return
	}

	if ctx.TargetEntity == nil {
		return
	}

	ctx.Game.Freeze(ctx.TargetEntity)
}
