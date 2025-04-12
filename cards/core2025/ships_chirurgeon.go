package cards

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

type ShipsChirurgeon struct {
	triggerId int
}

func (s *ShipsChirurgeon) Register(cm *game.CardManager) {
	card := game.Card{
		Name:        "Ship's Chirurgeon",
		ZhName:      "随船外科医师",
		ID:          "CORE_WON_065",
		Description: "在你召唤一个随从后，使其获得+1生命值。",
		Type:        game.Minion,
		Cost:        1,
		Health:      2,
		Attack:      1,
		Load:        s.Load,
		Unload:      s.Unload,
	}

	cm.RegisterCard(card)
}

func (s *ShipsChirurgeon) Load(g *game.Game, self *game.Entity) {
	s.triggerId = g.TriggerManager.RegisterTrigger(game.TriggerMinionSummoned, self, s.OnSummon, false)
}

func (s *ShipsChirurgeon) Unload(g *game.Game, self *game.Entity) {
	g.TriggerManager.UnregisterTrigger(s.triggerId)
}

func (s *ShipsChirurgeon) OnSummon(ctx *game.TriggerContext, self *game.Entity) {
	if self.CurrentZone != game.ZONE_PLAY { // only work when self is in play
		return
	}

	summonedMinion := ctx.SourceEntity
	if summonedMinion == nil {
		return
	}

	logger.Info("ShipsChirurgeon OnSummon", logger.String("summonedMinion", summonedMinion.Card.Name))

	if summonedMinion == self { // do not work on self
		return
	}
	if summonedMinion.Owner != self.Owner { // only work on friendly minion
		return
	}

	summonedMinion.Health++
	summonedMinion.MaxHealth++
}
