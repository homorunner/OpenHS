package cards

import (
	"github.com/openhs/internal/game"
)

type Malfurion struct{}

func (m *Malfurion) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Malfurion Stormrage",
		ZhName: "玛法里奥·怒风",
		ID:     "HERO_06",
		Health: 30,
		Type:   game.Hero,
	}

	cm.RegisterCard(card)
}
