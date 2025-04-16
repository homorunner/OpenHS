package game

import (
	"errors"

	"github.com/openhs/internal/logger"
)

// PlayWeapon handles playing a weapon card
func (g *Game) PlayWeapon(player *Player, entity *Entity, target *Entity) error {
	// If player already has a weapon, move it to graveyard
	if player.Weapon != nil {
		player.Graveyard = append(player.Graveyard, player.Weapon)
		player.Weapon.CurrentZone = ZONE_GRAVEYARD
		logger.Debug("Weapon moved to graveyard", logger.String("name", player.Weapon.Card.Name))
	}

	// Equip the new weapon
	player.Weapon = entity
	entity.CurrentZone = ZONE_PLAY
	logger.Debug("Weapon equipped", logger.String("name", entity.Card.Name))

	return nil
}

// DecreaseWeaponDurability decreases the durability of the weapon
// Note: this function will not destroy the weapon, that is handled in processGraveyard()
func (g *Game) DecreaseWeaponDurability(player *Player) error {
	if player.Weapon == nil {
		return errors.New("no weapon equipped")
	}

	weapon := player.Weapon
	weapon.Health--
	logger.Debug("Weapon durability decreased", logger.String("name", weapon.Card.Name))

	return nil
}
