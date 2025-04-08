package engine

import (
	"errors"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Attack handles combat between an attacker and defender entity
// It processes the damage exchange and any special effects
func (e *Engine) Attack(attacker *game.Entity, defender *game.Entity, skipValidation bool) error {
	// Validate the attack
	if !skipValidation {
		if err := e.validateAttack(attacker, defender); err != nil {
			return err
		}
	}

	// Set the phase to main combat
	e.game.Phase = game.MainCombat

	logger.Info("Attack initiated",
		logger.String("attacker", attacker.Card.Name),
		logger.String("defender", defender.Card.Name))

	// Process pre-attack triggers or effects if needed

	// Get attack values
	attackerDamage := attacker.Attack
	defenderDamage := defender.Attack

	// Decrease weapon durability if attacker is a hero
	var attackerWeapon *game.Entity
	if attacker.Card.Type == game.Hero && attacker.Owner != nil && attacker.Owner.Weapon != nil {
		attackerWeapon = attacker.Owner.Weapon
		e.decreaseWeaponDurability(attacker.Owner)
	}

	// Deal damage simultaneously
	// TODO: maybe merge the poisonous logic to DealDamage function
	if attackerDamage > 0 {
		e.DealDamage(attacker, defender, attackerDamage)

		// Check for poisonous effect on attacker
		if defender.Card.Type == game.Minion {
			if (attackerWeapon != nil && game.HasTag(attackerWeapon.Tags, game.TAG_POISONOUS)) ||
				game.HasTag(attacker.Tags, game.TAG_POISONOUS) {
				// If defender is still alive after taking damage, mark it for destruction
				defender.IsDestroyed = true
				logger.Info("Poisonous effect triggered",
					logger.String("source", attacker.Card.Name),
					logger.String("target", defender.Card.Name))
			}
		}
	}
	if defenderDamage > 0 {
		e.DealDamage(defender, attacker, defenderDamage)

		// Check for poisonous effect on defender
		if attacker.Card.Type == game.Minion && game.HasTag(defender.Tags, game.TAG_POISONOUS) {
			// If attacker is still alive after taking damage, mark it for destruction
			attacker.IsDestroyed = true
			logger.Info("Poisonous effect triggered",
				logger.String("source", defender.Card.Name),
				logger.String("target", attacker.Card.Name))
		}
	}

	// Mark attacker as having attacked this turn
	attacker.NumAttackThisTurn++

	// Check if attacker is exhausted
	expectedAttacks := 1
	if game.HasTag(attacker.Tags, game.TAG_WINDFURY) ||
		(attackerWeapon != nil && game.HasTag(attackerWeapon.Tags, game.TAG_WINDFURY)) {
		expectedAttacks = 2
	}
	if attacker.NumAttackThisTurn >= expectedAttacks {
		attacker.Exhausted = true
	}

	// Process special effects like poison, freeze, etc.

	// Process post-attack triggers or effects if needed

	// Check for deaths and update aura
	e.processDestroyAndUpdateAura()

	// Set the phase back to main
	e.game.Phase = game.MainAction

	return nil
}

// validateAttack checks if the attack is legal
func (e *Engine) validateAttack(attacker *game.Entity, defender *game.Entity) error {
	// Check for nil entities
	if attacker == nil || defender == nil {
		return errors.New("invalid attacker or defender")
	}

	// Check if attacker can attack
	if attacker.Attack <= 0 {
		return errors.New("attacker has 0 or negative attack")
	}

	// Check if attacker is exhausted
	if attacker.Exhausted {
		return errors.New("entity is exhausted and cannot attack this turn")
	}

	// Check if defender is a valid target
	// Only minion or hero of another player can be attack target
	if defender.Card.Type != game.Minion && defender.Card.Type != game.Hero {
		return errors.New("defender must be a minion or hero")
	}

	// Check if defender belongs to a different player
	if defender.Owner == attacker.Owner {
		return errors.New("cannot attack your own entities")
	}

	// Check rush restriction - Entities with rush cannot attack heroes on their first turn in field
	if game.HasTag(attacker.Tags, game.TAG_RUSH) &&
		defender.Card.Type == game.Hero &&
		attacker.NumTurnInPlay == 0 {
		return errors.New("minions with rush cannot attack heroes on their first turn")
	}

	// Additional validation logic can be added here
	// - Check for special effects like taunt, stealth, etc.

	return nil
}

func (e *Engine) processDestroyAndUpdateAura() {
	// Update aura

	// Trigger summon events

	// Process destroy, trigger deathrattle and reborn events (loop until no more entity dies)
	for e.processDestroyedWeapons() || e.processGraveyard() {
		e.processReborn()
	}

	// Update aura
}

func (e *Engine) processDestroyedWeapons() bool {
	destroyed := false
	for _, player := range e.game.Players {
		if player.Weapon != nil && (player.Weapon.Health <= 0 || player.Weapon.IsDestroyed) {
			player.Graveyard = append(player.Graveyard, player.Weapon)
			player.Weapon = nil
			destroyed = true
		}
	}
	return destroyed
}

func (e *Engine) processGraveyard() bool {
	destroyed := false
	for _, player := range e.game.Players {
		for _, minion := range player.Field {
			if minion.Health <= 0 || minion.IsDestroyed {
				destroyed = true
				e.removeEntityFromBoard(player, minion)

				// TODO: trigger death, deathrattle, infuse, add to reborn list, etc.

				// add to graveyard
				player.Graveyard = append(player.Graveyard, minion)
			}
		}
	}

	return destroyed
}

func (e *Engine) processReborn() {
}
