package game

import "github.com/openhs/internal/logger"

// AddEntityToField adds a new entity to a player's field at the specified position
// If fieldPos is -1, it will be added to the end
// If field is full, the entity is moved to ZONE_NONE and returns false
func (g *Game) AddEntityToField(player *Player, entity *Entity, fieldPos int) bool {
	// Check if field is full
	if len(player.Field) >= player.FieldSize {
		entity.CurrentZone = ZONE_NONE
		return false
	}

	if fieldPos == -1 {
		fieldPos = len(player.Field)
	}

	// Set the entity's zone to PLAY
	oldZone := entity.CurrentZone
	entity.CurrentZone = ZONE_PLAY

	// Reset attack and turn counters
	entity.NumAttackThisTurn = 0
	entity.NumTurnInPlay = 0 // First turn in play

	// Handle charge and rush tags
	if HasTag(entity.Tags, TAG_CHARGE) {
		// Minions with charge can attack immediately
		entity.Exhausted = false
		logger.Info("Charge effect activated", logger.String("name", entity.Card.Name))
	} else if HasTag(entity.Tags, TAG_RUSH) {
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
		player.Field = append(player.Field[:fieldPos], append([]*Entity{entity}, player.Field[fieldPos:]...)...)
	}

	// Trigger minion summoned event if coming from a different zone
	if oldZone != ZONE_PLAY {
		minionSummonedCtx := TriggerContext{
			Game:         g,
			SourceEntity: entity,
			Phase:        g.Phase,
		}
		g.TriggerManager.ActivateTrigger(TriggerMinionSummoned, minionSummonedCtx)
	}

	return true
}

// AddEntityToHand adds an entity to a player's hand
// handPos is the position in the hand (-1 for end)
// If hand is full, the entity is moved to ZONE_REMOVEDFROMGAME and returns false
func (g *Game) AddEntityToHand(player *Player, entity *Entity, handPos int) (*Entity, bool) {
	if handPos == -1 {
		handPos = len(player.Hand)
	}

	if len(player.Hand) >= player.HandSize {
		entity.CurrentZone = ZONE_REMOVEDFROMGAME
		return nil, false
	}

	// Set the entity's zone to HAND
	entity.CurrentZone = ZONE_HAND

	// Add entity to hand at the specified position
	if handPos == len(player.Hand) {
		player.Hand = append(player.Hand, entity)
	} else {
		player.Hand = append(player.Hand[:handPos], append([]*Entity{entity}, player.Hand[handPos:]...)...)
	}

	return entity, true
}

// MoveFromDeckToHand moves an entity from a player's deck to their hand
// deckIndex is the index in the deck
// handPos is the position in the hand (-1 for end)
// Returns true if the entity is moved to the hand, false if it is discarded
func (g *Game) MoveFromDeckToHand(player *Player, deckIndex, handPos int) (*Entity, bool) {
	// Get the entity from deck
	entity := player.Deck[deckIndex]

	// Remove from deck
	player.Deck = append(player.Deck[:deckIndex], player.Deck[deckIndex+1:]...)

	// Add to hand
	return g.AddEntityToHand(player, entity, handPos)
}

// MoveFromHandToField moves an entity from a player's hand to their field
// handIndex is the index in the hand
// fieldPos is the position on the field (-1 for end)
// Returns true if the entity is moved to the field, false if it is discarded
func (g *Game) MoveFromHandToField(player *Player, handIndex, fieldPos int) bool {
	// Get the entity from hand
	entity := player.Hand[handIndex]

	// Remove from hand
	player.Hand = append(player.Hand[:handIndex], player.Hand[handIndex+1:]...)

	// Add to field
	return g.AddEntityToField(player, entity, fieldPos)
}

// MoveFromDeckToField moves an entity from a player's deck to their field
// deckIndex is the index in the deck
// fieldPos is the position on the field (-1 for end)
// Returns true if the entity is moved to the field, false if it is discarded
func (g *Game) MoveFromDeckToField(player *Player, deckIndex, fieldPos int) bool {
	// Get the entity from deck
	entity := player.Deck[deckIndex]

	// Remove from deck
	player.Deck = append(player.Deck[:deckIndex], player.Deck[deckIndex+1:]...)

	// Add to field
	return g.AddEntityToField(player, entity, fieldPos)
}

// Helper to remove entity from board
func (g *Game) removeEntityFromBoard(player *Player, entity *Entity) {
	// Find the entity in the player's field and remove it
	for i, fieldEntity := range player.Field {
		if fieldEntity == entity {
			// Remove the entity by replacing it with the last element and then trimming the slice
			lastIdx := len(player.Field) - 1
			player.Field[i] = player.Field[lastIdx]
			player.Field = player.Field[:lastIdx]
			break
		}
	}
}
