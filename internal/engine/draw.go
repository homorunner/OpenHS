package engine

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// DrawCard draws a card from the player's deck to their hand
// It returns the drawn entity or nil if no card was drawn
func (e *Engine) DrawCard(player *game.Player) *game.Entity {
	return e.DrawSpecificCard(player, "")
}

// DrawSpecificCard allows drawing a specific card from the deck
// If cardToDraw is nil, it draws the top card of the deck
//
// Returns the drawn entity or nil if no card was drawn
func (e *Engine) DrawSpecificCard(player *game.Player, cardToDraw string) *game.Entity {
	// Check if the deck is empty
	if len(player.Deck) == 0 {
		if cardToDraw == "" {
			logger.Info("Player taking fatigue damage", logger.Int("damage", player.FatigueDamage))
			player.FatigueDamage++
			e.DealDamage(nil, player.Hero, player.FatigueDamage)
		}

		return nil
	}

	var drawIndex int

	// If a specific card is requested, find and draw it
	if cardToDraw != "" {
		found := false
		// Find the specific card in the deck
		for i, deckEntity := range player.Deck {
			if deckEntity.Card.Name == cardToDraw {
				// Get the entity and remove it from the deck
				drawIndex = i
				found = true
				break
			}
		}

		// If the card wasn't found, return nil
		if !found {
			logger.Info("Card not found in deck", logger.String("card", cardToDraw))
			return nil
		}
	} else {
		// Get the top entity from the deck
		drawIndex = len(player.Deck) - 1
	}

	// Try to add entity to hand
	entity, ok := e.MoveFromDeckToHand(player, drawIndex, -1)
	if !ok {
		return nil
	}

	// Create context for card drawn trigger
	cardDrawnCtx := game.TriggerContext{
		Game:         e.game,
		SourceEntity: entity,
		TargetEntity: player.Hero, // Associate with the player's hero
		Phase:        e.game.Phase,
	}

	// Trigger card drawn event
	e.game.TriggerManager.ActivateTrigger(game.TriggerCardDrawn, cardDrawnCtx)

	return entity
}
