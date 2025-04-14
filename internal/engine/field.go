package engine

import (
	"errors"

	"github.com/openhs/internal/game"
	"github.com/openhs/internal/logger"
)

// Field-related errors
var (
	ErrBattlefieldFull  = errors.New("battlefield is full")
	ErrInvalidHandIndex = errors.New("invalid hand index")
	ErrInvalidDeckIndex = errors.New("invalid deck index")
)

// AddEntityToField adds a new entity to a player's field at the specified position
// If fieldPos is -1, it will be added to the end
// Returns an error if the field is full
func (e *Engine) AddEntityToField(player *game.Player, entity *game.Entity, fieldPos int) error {
	// Check if field is full
	if len(player.Field) >= player.FieldSize {
		return ErrBattlefieldFull
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
	if fieldPos < 0 || fieldPos > len(player.Field) {
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

	return nil
}

// MoveFromHandToField moves an entity from a player's hand to their field
// handIndex is the index in the hand
// fieldPos is the position on the field (-1 for end)
// Returns an error if the field is full or the hand index is invalid
func (e *Engine) MoveFromHandToField(player *game.Player, handIndex, fieldPos int) error {
	// Validate hand index
	if handIndex < 0 || handIndex >= len(player.Hand) {
		return ErrInvalidHandIndex
	}

	// Check if field is full
	if len(player.Field) >= player.FieldSize {
		return ErrBattlefieldFull
	}

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
// Returns an error if the field is full or the deck index is invalid
func (e *Engine) MoveFromDeckToField(player *game.Player, deckIndex, fieldPos int) error {
	// Validate deck index
	if deckIndex < 0 || deckIndex >= len(player.Deck) {
		return ErrInvalidDeckIndex
	}

	// Check if field is full
	if len(player.Field) >= player.FieldSize {
		return ErrBattlefieldFull
	}

	// Get the entity from deck
	entity := player.Deck[deckIndex]

	// Remove from deck
	player.Deck = append(player.Deck[:deckIndex], player.Deck[deckIndex+1:]...)

	// Add to field
	return e.AddEntityToField(player, entity, fieldPos)
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
