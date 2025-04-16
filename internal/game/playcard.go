package game

import (
	"errors"

	"github.com/openhs/internal/logger"
)

// PlayCard processes playing a card from the player's hand
// Parameters:
// - player: The player playing the card
// - handIndex: The index of the card in the player's hand
// - target: Optional target for the card (can be nil)
// - fieldPos: Position on the field for minions (-1 for auto-positioning)
// - chooseOne: Index for choose one effects (0 for default choice)
func (g *Game) PlayCard(player *Player, handIndex int, target *Entity, fieldPos int, chooseOne int) error {
	// Validate hand index
	if handIndex < 0 || handIndex >= len(player.Hand) {
		return errors.New("invalid hand index")
	}

	// Get the entity to play
	entity := player.Hand[handIndex]

	// Check if we can play this card
	if err := g.TestPlayCard(player, entity, target, chooseOne); err != nil {
		return err
	}

	// Check field space for minions
	if entity.Card.Type == Minion && len(player.Field) >= player.FieldSize {
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
	entity.CurrentZone = ZONE_NONE

	// Process based on card type
	switch entity.Card.Type {
	case Minion:
		return g.PlayMinion(player, entity, target, fieldPos, chooseOne)
	case Spell:
		return g.PlaySpell(player, entity, target, chooseOne)
	case Weapon:
		return g.PlayWeapon(player, entity, target)
	case Hero:
		return g.PlayHero(player, entity, target, chooseOne)
	default:
		return errors.New("invalid card type")
	}
}

// TestPlayCard checks if a card can be played
func (g *Game) TestPlayCard(player *Player, entity *Entity, target *Entity, chooseOne int) error {
	// Basic checks
	if entity.Card.Cost > player.Mana {
		return errors.New("not enough mana")
	}

	// TODO: Check card-specific play requirements and target validity

	return nil
}

// PlayMinion handles playing a minion card
func (g *Game) PlayMinion(player *Player, entity *Entity, target *Entity, fieldPos int, chooseOne int) error {
	logger.Info("Minion played", logger.String("name", entity.Card.Name))

	// Trigger card played event
	cardPlayedCtx := TriggerContext{
		Game:         g,
		SourceEntity: entity,
		TargetEntity: target,
		Phase:        g.Phase,
		ExtraData:    map[string]interface{}{"card_type": Minion},
	}
	g.TriggerManager.ActivateTrigger(TriggerCardPlayed, cardPlayedCtx)

	// Try to add minion to the field at the specified position
	if g.AddEntityToField(player, entity, fieldPos) {
		// TODO: Process battlecry, triggers, etc.

		return nil
	}

	return nil
}

// PlaySpell handles playing a spell card
func (g *Game) PlaySpell(player *Player, entity *Entity, target *Entity, chooseOne int) error {
	// TODO: Add spell counters to Player struct

	logger.Info("Spell played", logger.String("name", entity.Card.Name))

	// Create context for card played trigger
	cardPlayedCtx := TriggerContext{
		Game:         g,
		SourceEntity: entity,
		TargetEntity: target,
		Phase:        g.Phase,
		ExtraData:    map[string]interface{}{"card_type": Spell},
	}

	// Trigger card played event
	g.TriggerManager.ActivateTrigger(TriggerCardPlayed, cardPlayedCtx)

	// Process spell effects
	for _, power := range entity.Card.Powers {
		if power.Type == PowerTypeSpell {
			power.Action(g, entity, target)
		}
	}
	g.processDestroyAndUpdateAura()

	// Move to graveyard after use
	player.Graveyard = append(player.Graveyard, entity)

	// Update the entity's zone
	entity.CurrentZone = ZONE_GRAVEYARD

	return nil
}

// PlayHero handles playing a hero card
func (g *Game) PlayHero(player *Player, entity *Entity, target *Entity, chooseOne int) error {
	// Store the old hero in graveyard
	if player.Hero != nil {
		player.Graveyard = append(player.Graveyard, player.Hero)
		// Update the old hero's zone
		player.Hero.CurrentZone = ZONE_GRAVEYARD
	}

	// Copy the old hero's state of attacking
	entity.Exhausted = player.Hero.Exhausted
	entity.NumAttackThisTurn = player.Hero.NumAttackThisTurn
	entity.NumTurnInPlay = 0 // First turn in play

	// Set the new hero
	player.Hero = entity

	// Update the entity's zone
	entity.CurrentZone = ZONE_PLAY

	logger.Info("Hero replaced", logger.String("name", entity.Card.Name))

	return nil
}
