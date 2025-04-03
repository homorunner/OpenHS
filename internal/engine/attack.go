package engine

import (
	"errors"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Attack handles combat between an attacker and defender card
// It processes the damage exchange and any special effects
func (e *Engine) Attack(attacker *game.Card, defender *game.Card, skipValidation bool) error {
	// Validate the attack
	if !skipValidation {
		if err := e.validateAttack(attacker, defender); err != nil {
			return err
		}
	}

	// Set the phase to main combat
	e.game.Phase = game.MainCombat

	logger.Info("Attack initiated", 
		logger.String("attacker", attacker.Name), 
		logger.String("defender", defender.Name))

	// Process pre-attack triggers or effects if needed

	// Get attack values
	attackerDamage := attacker.Attack
	defenderDamage := defender.Attack

	// Deal damage simultaneously
	if attackerDamage > 0 {
		if attacker.Type == game.Hero {
			// TODO: decrease weapon durability.
		}
		e.TakeDamage(defender, attackerDamage)
	}
	if defenderDamage > 0 {
		e.TakeDamage(attacker, defenderDamage)
	}

	// Process special effects like poison, freeze, etc.

	// Mark attacker as having attacked this turn

	// Process post-attack triggers or effects if needed

	// Check for deaths
	e.checkForDeaths(attacker, defender)

	// Set the phase back to main
	e.game.Phase = game.MainAction

	return nil
}

// validateAttack checks if the attack is legal
func (e *Engine) validateAttack(attacker *game.Card, defender *game.Card) error {
	// Check for nil cards
	if attacker == nil || defender == nil {
		return errors.New("invalid attacker or defender")
	}

	// Check if attacker can attack
	if attacker.Attack <= 0 {
		return errors.New("attacker has 0 or negative attack")
	}

	// Additional validation logic can be added here
	// - Check if attacker is exhausted
	// - Check if defender is a valid target
	// - Check for special effects like taunt, stealth, etc.

	return nil
}

// checkForDeaths checks if any cards have died after the attack
func (e *Engine) checkForDeaths(attacker *game.Card, defender *game.Card) {
	// Check if cards have died (health <= 0)
	if attacker.Health <= 0 {
		logger.Info("Card has died", logger.String("card", attacker.Name))
		// Process death effects
		// TODO: Implement death handling
	}

	if defender.Health <= 0 {
		logger.Info("Card has died", logger.String("card", defender.Name))
		// Process death effects
		// TODO: Implement death handling
	}
} 