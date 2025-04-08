package engine

import (
	"errors"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// PlayCard processes playing a card from the player's hand
// Parameters:
// - player: The player playing the card
// - handIndex: The index of the card in the player's hand
// - target: Optional target for the card (can be nil)
// - fieldPos: Position on the field for minions (-1 for auto-positioning)
// - chooseOne: Index for choose one effects (0 for default choice)
func (e *Engine) PlayCard(player *game.Player, handIndex int, target *game.Entity, fieldPos int, chooseOne int) error {
	// Validate hand index
	if handIndex < 0 || handIndex >= len(player.Hand) {
		return errors.New("invalid hand index")
	}

	// Get the entity to play
	entity := player.Hand[handIndex]

	// Check if we can play this card
	if err := e.testPlayCard(player, entity, target, chooseOne); err != nil {
		return err
	}

	// Check field space for minions
	if entity.Card.Type == game.Minion && len(player.Field) >= player.FieldSize {
		return errors.New("battlefield is full")
	}

	// Spend mana to play the card
	if entity.Card.Cost > 0 {
		// TODO: Handle temporary mana and special cost modifiers
		if player.Mana < entity.Card.Cost {
			return errors.New("not enough mana")
		}
		player.Mana -= entity.Card.Cost
	}

	// Record play history and update game state
	// TODO: Add NumCardsPlayedThisTurn to Player struct

	// Remove entity from hand
	player.Hand = append(player.Hand[:handIndex], player.Hand[handIndex+1:]...)

	// Process based on card type
	switch entity.Card.Type {
	case game.Minion:
		return e.playMinion(player, entity, target, fieldPos, chooseOne)
	case game.Spell:
		return e.playSpell(player, entity, target, chooseOne)
	case game.Weapon:
		return e.playWeapon(player, entity, target)
	case game.Hero:
		return e.playHero(player, entity, target, chooseOne)
	default:
		return errors.New("invalid card type")
	}
}

// canPlayCard checks if a card can be played
func (e *Engine) testPlayCard(player *game.Player, entity *game.Entity, target *game.Entity, chooseOne int) error {
	// Basic checks
	if entity.Card.Cost > player.Mana {
		return errors.New("not enough mana")
	}

	// TODO: Check card-specific play requirements and target validity

	return nil
}

// playMinion handles playing a minion card
func (e *Engine) playMinion(player *game.Player, entity *game.Entity, target *game.Entity, fieldPos int, chooseOne int) error {
	// The minion is exhausted when it enters play, unless it has charge
	// Note: exhausted and sleeping are not the same state, but for now we use the same flag
	entity.NumAttackThisTurn = 0

	if game.HasTag(entity.Tags, game.TAG_CHARGE) {
		// Minions with charge can attack immediately
		entity.Exhausted = false
		logger.Info("Charge effect activated", logger.String("name", entity.Card.Name))
	} else {
		entity.Exhausted = true
	}

	// Add minion to the field at the specified position
	if fieldPos < 0 || fieldPos > len(player.Field) {
		// Auto-position at the end
		player.Field = append(player.Field, entity)
	} else {
		// Insert at specified position
		player.Field = append(player.Field[:fieldPos], append([]*game.Entity{entity}, player.Field[fieldPos:]...)...)
	}

	logger.Info("Minion played", logger.String("name", entity.Card.Name))

	// TODO: Process battlecry, triggers, etc.

	return nil
}

// playSpell handles playing a spell card
func (e *Engine) playSpell(player *game.Player, entity *game.Entity, target *game.Entity, chooseOne int) error {
	// TODO: Add spell counters to Player struct

	logger.Info("Spell played", logger.String("name", entity.Card.Name))

	// TODO: Process spell effects

	// Move to graveyard after use
	player.Graveyard = append(player.Graveyard, entity)

	return nil
}

// playHero handles playing a hero card
func (e *Engine) playHero(player *game.Player, entity *game.Entity, target *game.Entity, chooseOne int) error {
	// Store the old hero in graveyard
	if player.Hero != nil {
		player.Graveyard = append(player.Graveyard, player.Hero)
	}

	// Copy the old hero's state of attacking
	entity.Exhausted = player.Hero.Exhausted
	entity.NumAttackThisTurn = player.Hero.NumAttackThisTurn

	// Set the new hero
	player.Hero = entity

	logger.Info("Hero replaced", logger.String("name", entity.Card.Name))

	return nil
}
