package game

import (
	"errors"

	"github.com/openhs/internal/logger"
)

// Attack handles combat between an attacker and defender entity
// It processes the damage exchange and any special effects
func (g *Game) Attack(attacker *Entity, defender *Entity, skipValidation bool) error {
	// Validate the attack
	if !skipValidation {
		if err := g.validateAttack(attacker, defender); err != nil {
			return err
		}
	}

	// Set the phase to main combat
	g.Phase = MainCombat

	logger.Info("Attack initiated",
		logger.String("attacker", attacker.Card.Name),
		logger.String("defender", defender.Card.Name))

	// TODO: Process pre-attack triggers

	// Get attack values
	attackerDamage := attacker.Attack
	defenderDamage := defender.Attack

	// Decrease weapon durability if attacker is a hero
	var attackerWeapon *Entity
	if attacker.Card.Type == Hero && attacker.Owner != nil && attacker.Owner.Weapon != nil {
		attackerWeapon = attacker.Owner.Weapon
		g.decreaseWeaponDurability(attacker.Owner)
	}

	// Deal damage simultaneously
	if attackerDamage > 0 {
		g.DealDamage(attacker, defender, attackerDamage)

		// Check for poisonous effect on attacker
		if defender.Card.Type == Minion {
			if (attackerWeapon != nil && HasTag(attackerWeapon.Tags, TAG_POISONOUS)) ||
				HasTag(attacker.Tags, TAG_POISONOUS) {
				// If defender is still alive after taking damage, mark it for destruction
				defender.IsDestroyed = true
				logger.Info("Poisonous effect triggered",
					logger.String("source", attacker.Card.Name),
					logger.String("target", defender.Card.Name))
			}
		}
	}
	if defenderDamage > 0 {
		g.DealDamage(defender, attacker, defenderDamage)

		// Check for poisonous effect on defender
		if attacker.Card.Type == Minion && HasTag(defender.Tags, TAG_POISONOUS) {
			// If attacker is still alive after taking damage, mark it for destruction
			attacker.IsDestroyed = true
			logger.Info("Poisonous effect triggered",
				logger.String("source", defender.Card.Name),
				logger.String("target", attacker.Card.Name))
		}
	}

	// Mark attacker as having attacked this turn
	attacker.NumAttackThisTurn++

	// Check if attacker is exhausted
	expectedAttacks := 1
	if HasTag(attacker.Tags, TAG_WINDFURY) ||
		(attackerWeapon != nil && HasTag(attackerWeapon.Tags, TAG_WINDFURY)) {
		expectedAttacks = 2
	}
	if attacker.NumAttackThisTurn >= expectedAttacks {
		attacker.Exhausted = true
	}

	// Process special effects like poison, freeze, etc.

	// Process post-attack triggers or effects if needed

	// Check for deaths and update aura
	g.processDestroyAndUpdateAura()

	// Set the phase back to main
	g.Phase = MainAction

	return nil
}

// validateAttack checks if the attack is legal
func (g *Game) validateAttack(attacker *Entity, defender *Entity) error {
	// Check for nil entities
	if attacker == nil || defender == nil {
		return errors.New("invalid attacker or defender")
	}

	// Check if attacker can attack
	if attacker.Attack <= 0 {
		return errors.New("attacker has 0 or negative attack")
	}

	// Check if attacker is exhausted
	if attacker.Exhausted {
		return errors.New("entity is exhausted and cannot attack this turn")
	}

	// Check if attacker is Frozen
	if HasTag(attacker.Tags, TAG_FROZEN) {
		return errors.New("frozen entities cannot attack")
	}

	// Check if defender is a valid target
	// Only minion or hero of another player can be attack target
	if defender.Card.Type != Minion && defender.Card.Type != Hero {
		return errors.New("defender must be a minion or hero")
	}

	// Check if defender belongs to a different player
	if defender.Owner == attacker.Owner {
		return errors.New("cannot attack your own entities")
	}

	// Check rush restriction - Entities with rush cannot attack heroes on their first turn in field
	if HasTag(attacker.Tags, TAG_RUSH) &&
		defender.Card.Type == Hero &&
		attacker.NumTurnInPlay == 0 {
		return errors.New("minions with rush cannot attack heroes on their first turn")
	}

	// Additional validation logic can be added here
	// - Check for special effects like taunt, stealth, etc.

	return nil
}

func (g *Game) processDestroyAndUpdateAura() {
	// Update aura

	// Trigger summon events

	// Process destroy, trigger deathrattle and reborn events (loop until no more entity dies)
	for g.processDestroyedWeapons() || g.ProcessGraveyard() {
		g.processReborn()
	}

	// Update aura
}

func (g *Game) processDestroyedWeapons() bool {
	destroyed := false
	for _, player := range g.Players {
		if player.Weapon != nil && (player.Weapon.Health <= 0 || player.Weapon.IsDestroyed) {
			player.Graveyard = append(player.Graveyard, player.Weapon)
			player.Weapon.CurrentZone = ZONE_GRAVEYARD
			player.Weapon = nil
			destroyed = true
		}
	}
	return destroyed
}

func (g *Game) ProcessGraveyard() bool {
	destroyed := false
	for _, player := range g.Players {
		for _, minion := range player.Field {
			if minion.Health <= 0 || minion.IsDestroyed {
				destroyed = true
				g.removeEntityFromBoard(player, minion)

				// Create context for minion death trigger
				deathCtx := TriggerContext{
					Game:         g,
					SourceEntity: minion,
					Phase:        g.Phase,
				}

				// Trigger minion death event
				g.TriggerManager.ActivateTrigger(TriggerMinionDeath, deathCtx)

				// TODO: trigger death, deathrattle, infuse, add to reborn list, etc.

				// add to graveyard
				player.Graveyard = append(player.Graveyard, minion)
				minion.CurrentZone = ZONE_GRAVEYARD
			}
		}
	}

	return destroyed
}

func (g *Game) processReborn() {
	// Reborn implementation
}

// ProcessAttack handles an attack from one entity to another
func (g *Game) ProcessAttack(attacker, defender *Entity) error {
	// Create context for attack triggers
	beforeAttackCtx := TriggerContext{
		Game:         g,
		SourceEntity: attacker,
		TargetEntity: defender,
		Phase:        g.Phase,
	}

	// Activate before attack triggers
	g.TriggerManager.ActivateTrigger(TriggerBeforeAttack, beforeAttackCtx)

	// Use the existing Attack method with skipValidation=false
	err := g.Attack(attacker, defender, false)
	if err != nil {
		return err
	}

	// Create context for after attack triggers
	afterAttackCtx := TriggerContext{
		Game:         g,
		SourceEntity: attacker,
		TargetEntity: defender,
		Phase:        g.Phase,
	}

	// Activate after attack triggers
	g.TriggerManager.ActivateTrigger(TriggerAfterAttack, afterAttackCtx)

	return nil
}

// CanAttack checks if an entity can attack a target
func (g *Game) CanAttack(attacker, defender *Entity) error {
	// Reuse the existing validation logic
	return g.validateAttack(attacker, defender)
}

// CanBeAttacked checks if an entity can be attacked
func (g *Game) CanBeAttacked(defender *Entity) error {
	// Check if defender is a valid target
	if defender == nil {
		return errors.New("invalid defender")
	}

	// Check if defender is a valid target type
	if defender.Card.Type != Minion && defender.Card.Type != Hero {
		return errors.New("defender must be a minion or hero")
	}

	// Check for special tags that prevent attacking
	if HasTag(defender.Tags, TAG_STEALTH) {
		return errors.New("stealthed entities cannot be attacked")
	}

	// Check for taunt restriction - if opponent has taunt, must attack that
	if !HasTag(defender.Tags, TAG_TAUNT) {
		opponentHasTaunt := false
		for _, player := range g.Players {
			if player == defender.Owner {
				continue
			}

			for _, minion := range player.Field {
				if HasTag(minion.Tags, TAG_TAUNT) {
					opponentHasTaunt = true
					break
				}
			}
		}

		if opponentHasTaunt {
			return errors.New("must attack entities with taunt first")
		}
	}

	return nil
}

// ShouldExhaustAfterAttack determines if an entity should be exhausted after attacking
func (g *Game) ShouldExhaustAfterAttack(entity *Entity) bool {
	if entity == nil {
		return false
	}

	// Calculate expected number of attacks based on entity's attributes
	expectedAttacks := 1
	if HasTag(entity.Tags, TAG_WINDFURY) {
		expectedAttacks = 2
	}

	return entity.NumAttackThisTurn >= expectedAttacks
}

// Helper to decrease weapon durability
func (g *Game) decreaseWeaponDurability(player *Player) {
	if player == nil || player.Weapon == nil {
		return
	}

	player.Weapon.Health--

	// Log the durability decrease
	logger.Debug("Weapon durability decreased",
		logger.String("weapon", player.Weapon.Card.Name),
		logger.Int("remaining", player.Weapon.Health))
}

// Freeze marks an entity as frozen. Frozen entities miss their next possible attack.
func (g *Game) Freeze(target *Entity) {
	// Only apply freeze if target is not already frozen
	if !HasTag(target.Tags, TAG_FROZEN) {
		target.Tags = append(target.Tags, NewTag(TAG_FROZEN, true))
		logger.Info("Entity frozen",
			logger.String("target", target.Card.Name))
	}
}
