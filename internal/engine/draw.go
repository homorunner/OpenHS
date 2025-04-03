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
		player.FatigueDamage++

		logger.Info("Player taking fatigue damage", logger.Int("damage", player.FatigueDamage))
		e.TakeDamage(player.Hero, player.FatigueDamage)
		return nil
	}
	
	var entity *game.Entity
	
	// If a specific card is requested, find and draw it
	if cardToDraw != "" {
		found := false
		// Find the specific card in the deck
		for i, deckEntity := range player.Deck {
			if deckEntity.Card.Name == cardToDraw {
				// Get the entity and remove it from the deck
				entity = player.Deck[i]
				player.Deck = append(player.Deck[:i], player.Deck[i+1:]...)
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
		entity = player.Deck[len(player.Deck)-1]
		player.Deck = player.Deck[:len(player.Deck)-1]
	}
	
	// Add entity to hand if space is available
	if !e.AddCardToHand(player, entity) {
		return nil
	}
	
	return entity
}

// AddCardToHand adds a card to the player's hand
// Returns true if successful, false if hand is full
func (e *Engine) AddCardToHand(player *game.Player, entity *game.Entity) bool {
	if len(player.Hand) >= player.HandSize {
		logger.Info("Hand is full, card discarded", logger.String("card", entity.Card.Name))
		return false
	}
	
	// Add entity to hand
	player.Hand = append(player.Hand, entity)
	
	return true
}
