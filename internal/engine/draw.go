package engine

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// DrawCard draws a card from the player's deck to their hand
// It returns the drawn card or nil if no card was drawn
func (e *Engine) DrawCard(player *game.Player) *game.Card {
	return e.DrawSpecificCard(player, "")
}

// DrawSpecificCard allows drawing a specific card from the deck
// If cardToDraw is nil, it draws the top card of the deck
// 
// Returns the drawn card or nil if no card was drawn
func (e *Engine) DrawSpecificCard(player *game.Player, cardToDraw string) *game.Card {
	// Check if the deck is empty
	if len(player.Deck) == 0 {
		player.FatigueDamage++

		logger.Info("Player taking fatigue damage", logger.Int("damage", player.FatigueDamage))
		e.TakeDamage(&player.Hero, player.FatigueDamage)
		return nil
	}
	
	var card game.Card
	
	// If a specific card is requested, find and draw it
	if cardToDraw != "" {
		found := false
		// Find the specific card in the deck
		for i, deckCard := range player.Deck {
			if deckCard.Name == cardToDraw {
				// Get the card and remove it from the deck
				card = player.Deck[i]
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
		// Get the top card from the deck
		card = player.Deck[len(player.Deck)-1]
		player.Deck = player.Deck[:len(player.Deck)-1]
	}
	
	// Add card to hand if space is available
	if !e.AddCardToHand(player, &card) {
		return nil
	}
	
	return &card
}

// AddCardToHand adds a card to the player's hand
// Returns true if successful, false if hand is full
func (e *Engine) AddCardToHand(player *game.Player, card *game.Card) bool {
	if len(player.Hand) >= player.HandSize {
		logger.Info("Hand is full, card discarded", logger.String("card", card.Name))
		return false
	}
	
	// Add card to hand
	player.Hand = append(player.Hand, *card)
	
	return true
}
