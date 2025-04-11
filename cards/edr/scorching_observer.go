package cards

import (
	"github.com/openhs/internal/game"
)

type ScorchingObserver struct{}

func (s *ScorchingObserver) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Scorching Observer",
		ZhName: "纵火眼魔",
		Cost:   9,
		Attack: 7,
		Health: 9,
		Type:   game.Minion,
		Tags: []game.Tag{
			game.NewTag(game.TAG_LIFESTEAL, true),
			game.NewTag(game.TAG_RUSH, true),
		},
	}

	cm.RegisterCard(card)
}
