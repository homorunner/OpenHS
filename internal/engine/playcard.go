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
func (e *Engine) PlayCard(player *game.Player, handIndex int, target *game.Card, fieldPos int, chooseOne int) error {
	// Validate hand index
	if handIndex < 0 || handIndex >= len(player.Hand) {
		return errors.New("invalid hand index")
	}

	// Get the card to play
	card := player.Hand[handIndex]

	// Check if we can play this card
	if err := e.testPlayCard(player, &card, target, chooseOne); err != nil {
		return err
	}

	// Check field space for minions
	if card.Type == game.Minion && len(player.Field) >= player.HandSize {
		return errors.New("battlefield is full")
	}

	// Spend mana to play the card
	if card.Cost > 0 {
		// TODO: Handle temporary mana and special cost modifiers
		if player.Mana < card.Cost {
			return errors.New("not enough mana")
		}
		player.Mana -= card.Cost
	}

	// Record play history and update game state
	// TODO: Add NumCardsPlayedThisTurn to Player struct
	
	// Remove card from hand
	player.Hand = append(player.Hand[:handIndex], player.Hand[handIndex+1:]...)

	// Process based on card type
	switch card.Type {
	case game.Minion:
		return e.playMinion(player, &card, target, fieldPos, chooseOne)
	case game.Spell:
		return e.playSpell(player, &card, target, chooseOne)
	case game.Weapon:
		return e.playWeapon(player, &card, target)
	case game.Hero:
		return e.playHero(player, &card, target, chooseOne)
	default:
		return errors.New("invalid card type")
	}
}

// canPlayCard checks if a card can be played
func (e *Engine) testPlayCard(player *game.Player, card *game.Card, target *game.Card, chooseOne int) error {
	// Basic checks
	if card.Cost > player.Mana {
		return errors.New("not enough mana")
	}

	// TODO: Check card-specific play requirements and target validity

	return nil
}

// playMinion handles playing a minion card
func (e *Engine) playMinion(player *game.Player, card *game.Card, target *game.Card, fieldPos int, chooseOne int) error {
	// Create a new minion entity (using Card for now since Minion type doesn't exist yet)
	// TODO: Create a proper Minion type
	minion := *card
	
	// Add minion to the field at the specified position
	if fieldPos < 0 || fieldPos > len(player.Field) {
		// Auto-position at the end
		player.Field = append(player.Field, minion)
	} else {
		// Insert at specified position
		player.Field = append(player.Field[:fieldPos], append([]game.Card{minion}, player.Field[fieldPos:]...)...)
	}

	logger.Info("Minion played", logger.String("name", card.Name))

	// TODO: Process battlecry, triggers, etc.
	
	return nil
}

// playSpell handles playing a spell card
func (e *Engine) playSpell(player *game.Player, card *game.Card, target *game.Card, chooseOne int) error {
	// TODO: Add spell counters to Player struct
	
	logger.Info("Spell played", logger.String("name", card.Name))

	// TODO: Process spell effects
	
	// Move to graveyard after use
	player.Graveyard = append(player.Graveyard, *card)
	
	return nil
}

// playWeapon handles playing a weapon card
func (e *Engine) playWeapon(player *game.Player, card *game.Card, target *game.Card) error {
	// TODO: Create a proper Weapon type
	weapon := *card

	// TODO: Implement weapon handling
	// For now, just replace the weapon
	if player.HasWeapon {
		player.Graveyard = append(player.Graveyard, player.Weapon)
	}
	player.Weapon = weapon
	player.HasWeapon = true

	logger.Info("Weapon equipped", logger.String("name", card.Name))
	
	return nil
}

// playHero handles playing a hero card
func (e *Engine) playHero(player *game.Player, card *game.Card, target *game.Card, chooseOne int) error {
	// TODO: Implement hero card handling
	// For now, just replace the hero
	player.Hero = *card

	logger.Info("Hero replaced", logger.String("name", card.Name))
	
	return nil
}

