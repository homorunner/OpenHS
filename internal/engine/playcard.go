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

	// Update entity zone (temporarily set to NONE while in transition)
	entity.CurrentZone = game.ZONE_NONE

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

// testPlayCard checks if a card can be played
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
	logger.Info("Minion played", logger.String("name", entity.Card.Name))

	// Trigger card played event
	cardPlayedCtx := game.TriggerContext{
		Game:         e.game,
		SourceEntity: entity,
		TargetEntity: target,
		Phase:        e.game.Phase,
		ExtraData:    map[string]interface{}{"card_type": game.Minion},
	}
	e.game.TriggerManager.ActivateTrigger(game.TriggerCardPlayed, cardPlayedCtx)

	// Add minion to the field at the specified position
	if err := e.AddEntityToField(player, entity, fieldPos); err != nil {
		return err
	}

	// TODO: Process battlecry, triggers, etc.

	return nil
}

// playSpell handles playing a spell card
func (e *Engine) playSpell(player *game.Player, entity *game.Entity, target *game.Entity, chooseOne int) error {
	// TODO: Add spell counters to Player struct

	logger.Info("Spell played", logger.String("name", entity.Card.Name))

	// Create context for card played trigger
	cardPlayedCtx := game.TriggerContext{
		Game:         e.game,
		SourceEntity: entity,
		TargetEntity: target,
		Phase:        e.game.Phase,
		ExtraData:    map[string]interface{}{"card_type": game.Spell},
	}

	// Trigger card played event
	e.game.TriggerManager.ActivateTrigger(game.TriggerCardPlayed, cardPlayedCtx)

	// TODO: Process spell effects

	// Move to graveyard after use
	player.Graveyard = append(player.Graveyard, entity)

	// Update the entity's zone
	entity.CurrentZone = game.ZONE_GRAVEYARD

	return nil
}

// playHero handles playing a hero card
func (e *Engine) playHero(player *game.Player, entity *game.Entity, target *game.Entity, chooseOne int) error {
	// Store the old hero in graveyard
	if player.Hero != nil {
		player.Graveyard = append(player.Graveyard, player.Hero)
		// Update the old hero's zone
		player.Hero.CurrentZone = game.ZONE_GRAVEYARD
	}

	// Copy the old hero's state of attacking
	entity.Exhausted = player.Hero.Exhausted
	entity.NumAttackThisTurn = player.Hero.NumAttackThisTurn
	entity.NumTurnInPlay = 0 // First turn in play

	// Set the new hero
	player.Hero = entity

	// Update the entity's zone
	entity.CurrentZone = game.ZONE_PLAY

	logger.Info("Hero replaced", logger.String("name", entity.Card.Name))

	return nil
}

// playWeapon is implemented in weapon.go
