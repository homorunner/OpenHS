package cards

import (
	"github.com/openhs/internal/game"
)

type WaterElemental struct{}

func (w *WaterElemental) Register(cm *game.CardManager) {
	card := game.Card{
		Name:   "Water Elemental",
		ZhName: "水元素",
		Cost:   4,
		Attack: 3,
		Health: 6,
		Type:   game.Minion,
	}

	cm.RegisterCard(card)
}
