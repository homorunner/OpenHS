package engine

import (
	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// AddEntityToField adds a new entity to a player's field at the specified position
// If fieldPos is -1, it will be added to the end
// If field is full, the entity is moved to ZONE_NONE and returns false
func (e *Engine) AddEntityToField(player *game.Player, entity *game.Entity, fieldPos int) bool {
	// Check if field is full
	if len(player.Field) >= player.FieldSize {
		entity.CurrentZone = game.ZONE_NONE
		return false
	}

	if fieldPos == -1 {
		fieldPos = len(player.Field)
	}

	// Set the entity's zone to PLAY
	oldZone := entity.CurrentZone
	entity.CurrentZone = game.ZONE_PLAY

	// Reset attack and turn counters
	entity.NumAttackThisTurn = 0
	entity.NumTurnInPlay = 0 // First turn in play

	// Handle charge and rush tags
	if game.HasTag(entity.Tags, game.TAG_CHARGE) {
		// Minions with charge can attack immediately
		entity.Exhausted = false
		logger.Info("Charge effect activated", logger.String("name", entity.Card.Name))
	} else if game.HasTag(entity.Tags, game.TAG_RUSH) {
		// Minions with rush can attack minions immediately, but not heroes
		entity.Exhausted = false
		logger.Info("Rush effect activated", logger.String("name", entity.Card.Name))
	} else {
		entity.Exhausted = true
	}

	// Add entity to field at the specified position
	if fieldPos == len(player.Field) {
		// Auto-position at the end
		player.Field = append(player.Field, entity)
	} else {
		// Insert at specified position
		player.Field = append(player.Field[:fieldPos], append([]*game.Entity{entity}, player.Field[fieldPos:]...)...)
	}

	// Trigger minion summoned event if coming from a different zone
	if oldZone != game.ZONE_PLAY {
		minionSummonedCtx := game.TriggerContext{
			Game:         e.game,
			SourceEntity: entity,
			Phase:        e.game.Phase,
		}
		e.game.TriggerManager.ActivateTrigger(game.TriggerMinionSummoned, minionSummonedCtx)
	}

	return true
}

// AddEntityToHand adds an entity to player's hand at the specified position
// If handPos is -1, it will be added to the end
// If hand is full, the entity is moved to ZONE_REMOVEDFROMGAME and returns false
func (e *Engine) AddEntityToHand(player *game.Player, entity *game.Entity, handPos int) (*game.Entity, bool) {
	if handPos == -1 {
		handPos = len(player.Hand)
	}

	if len(player.Hand) >= player.HandSize {
		entity.CurrentZone = game.ZONE_REMOVEDFROMGAME
		return nil, false
	}

	// Set the entity's zone to HAND
	entity.CurrentZone = game.ZONE_HAND

	// Add entity to hand at the specified position
	if handPos == len(player.Hand) {
		player.Hand = append(player.Hand, entity)
	} else {
		player.Hand = append(player.Hand[:handPos], append([]*game.Entity{entity}, player.Hand[handPos:]...)...)
	}

	return entity, true
}

// MoveFromHandToField moves an entity from a player's hand to their field
// handIndex is the index in the hand
// fieldPos is the position on the field (-1 for end)
// Returns true if the entity is moved to the field, false if it is discarded
func (e *Engine) MoveFromHandToField(player *game.Player, handIndex, fieldPos int) bool {
	// Get the entity from hand
	entity := player.Hand[handIndex]

	// Remove from hand
	player.Hand = append(player.Hand[:handIndex], player.Hand[handIndex+1:]...)

	// Add to field
	return e.AddEntityToField(player, entity, fieldPos)
}

// MoveFromDeckToField moves an entity from a player's deck to their field
// deckIndex is the index in the deck
// fieldPos is the position on the field (-1 for end)
// Returns true if the entity is moved to the field, false if it is discarded
func (e *Engine) MoveFromDeckToField(player *game.Player, deckIndex, fieldPos int) bool {
	// Get the entity from deck
	entity := player.Deck[deckIndex]

	// Remove from deck
	player.Deck = append(player.Deck[:deckIndex], player.Deck[deckIndex+1:]...)

	// Add to field
	return e.AddEntityToField(player, entity, fieldPos)
}

// MoveFromDeckToHand moves an entity from a player's deck to their hand
// deckIndex is the index in the deck
// handPos is the position in the hand (-1 for end)
// Returns true if the entity is moved to the hand, false if it is discarded
func (e *Engine) MoveFromDeckToHand(player *game.Player, deckIndex, handPos int) (*game.Entity, bool) {
	// Get the entity from deck
	entity := player.Deck[deckIndex]

	// Remove from deck
	player.Deck = append(player.Deck[:deckIndex], player.Deck[deckIndex+1:]...)

	// Add to hand
	return e.AddEntityToHand(player, entity, handPos)
}

func (e *Engine) removeEntityFromBoard(player *game.Player, entity *game.Entity) {
	for i, minion := range player.Field {
		if minion == entity {
			player.Field = append(player.Field[:i], player.Field[i+1:]...)
			// The entity is being removed from the board,
			// but it will be placed in a new zone afterwards,
			// so we set it to NONE temporarily
			entity.CurrentZone = game.ZONE_NONE
			break
		}
	}
}
