package engine

import (
	"errors"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// playWeapon handles playing a weapon card
func (e *Engine) playWeapon(player *game.Player, entity *game.Entity, target *game.Entity) error {
	// If player already has a weapon, move it to graveyard
	if player.Weapon != nil {
		player.Graveyard = append(player.Graveyard, player.Weapon)
		logger.Debug("Weapon moved to graveyard", logger.String("name", player.Weapon.Card.Name))
	}

	// Equip the new weapon
	player.Weapon = entity
	logger.Debug("Weapon equipped", logger.String("name", entity.Card.Name))
	
	return nil
}

func (e *Engine) decreaseWeaponDurability(player *game.Player) error {
	if player.Weapon == nil {
		return errors.New("no weapon equipped")
	}

	weapon := player.Weapon
	weapon.Health--
	logger.Debug("Weapon durability decreased", logger.String("name", weapon.Card.Name))

	if weapon.Health <= 0 {
		player.Graveyard = append(player.Graveyard, weapon)
		player.Weapon = nil
		logger.Debug("Weapon destroyed", logger.String("name", weapon.Card.Name))
	}

	return nil
}
