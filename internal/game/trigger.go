package game

import "github.com/openhs/internal/logger"

// TriggerType represents the type of trigger
type TriggerType int

const (
	// Turn triggers
	TriggerTurnStart TriggerType = iota
	TriggerTurnEnd

	// Card triggers
	TriggerCardPlayed
	TriggerCardDrawn

	// Combat triggers
	TriggerBeforeAttack
	TriggerAfterAttack
	TriggerDamageTaken
	TriggerHealReceived

	// Minion triggers
	TriggerMinionSummoned
	TriggerMinionDeath

	// Hero triggers
	TriggerHeroDamageTaken
	TriggerHeroPowerUsed

	// More triggers can be added here

	// Just for tracking the amount of triggers
	LastTrigger
)

func (t TriggerType) String() string {
	switch t {
	case TriggerTurnStart:
		return "TriggerTurnStart"
	case TriggerTurnEnd:
		return "TriggerTurnEnd"
	case TriggerCardPlayed:
		return "TriggerCardPlayed"
	case TriggerCardDrawn:
		return "TriggerCardDrawn"
	case TriggerBeforeAttack:
		return "TriggerBeforeAttack"
	case TriggerAfterAttack:
		return "TriggerAfterAttack"
	case TriggerDamageTaken:
		return "TriggerDamageTaken"
	case TriggerHealReceived:
		return "TriggerHealReceived"
	case TriggerMinionSummoned:
		return "TriggerMinionSummoned"
	case TriggerMinionDeath:
		return "TriggerMinionDeath"
	case TriggerHeroDamageTaken:
		return "TriggerHeroDamageTaken"
	case TriggerHeroPowerUsed:
		return "TriggerHeroPowerUsed"
	default:
		return "UnknownTrigger"
	}
}

// TriggerContext holds contextual information about a triggered event
type TriggerContext struct {
	Game         *Game
	SourceEntity *Entity                // Entity that caused the trigger
	TargetEntity *Entity                // Entity that was targeted (if applicable)
	Value        int                    // Generic value field for damage, healing, etc.
	Phase        GamePhase              // Current game phase
	ExtraData    map[string]interface{} // Additional data specific to the trigger
}

// TriggerFunc is a function called when a trigger activates
type TriggerFunc func(ctx *TriggerContext, self *Entity)

// TriggerRegistration represents a registered trigger
type TriggerRegistration struct {
	ID           int
	Type         TriggerType
	RegisteredOn *Entity // Entity that registered the trigger
	Callback     TriggerFunc
	OneTimeOnly  bool
}

// TriggerManager manages all game triggers
type TriggerManager struct {
	registrations map[TriggerType][]TriggerRegistration
	nextID        int
}

// NewTriggerManager creates a new trigger manager
func NewTriggerManager() *TriggerManager {
	manager := &TriggerManager{
		registrations: make(map[TriggerType][]TriggerRegistration),
		nextID:        1,
	}

	// Initialize slices for each trigger type
	for i := 0; i < int(LastTrigger); i++ {
		manager.registrations[TriggerType(i)] = make([]TriggerRegistration, 0)
	}

	return manager
}

// RegisterTrigger registers a new trigger function
func (tm *TriggerManager) RegisterTrigger(
	triggerType TriggerType,
	source *Entity,
	callback TriggerFunc,
	oneTimeOnly bool,
) int {
	registration := TriggerRegistration{
		ID:           tm.nextID,
		Type:         triggerType,
		RegisteredOn: source,
		Callback:     callback,
		OneTimeOnly:  oneTimeOnly,
	}

	tm.registrations[triggerType] = append(tm.registrations[triggerType], registration)
	tm.nextID++

	return registration.ID
}

// UnregisterTrigger removes a trigger by its ID
func (tm *TriggerManager) UnregisterTrigger(id int) bool {
	for triggerType, registrations := range tm.registrations {
		for i, reg := range registrations {
			if reg.ID == id {
				// Remove the registration by swapping with the last element and truncating
				lastIndex := len(registrations) - 1
				registrations[i] = registrations[lastIndex]
				tm.registrations[triggerType] = registrations[:lastIndex]
				return true
			}
		}
	}
	return false
}

// UnregisterAllForEntity removes all triggers for a specific entity
func (tm *TriggerManager) UnregisterAllForEntity(entity *Entity) {
	for triggerType, registrations := range tm.registrations {
		newRegistrations := make([]TriggerRegistration, 0, len(registrations))

		for _, reg := range registrations {
			if reg.RegisteredOn != entity {
				newRegistrations = append(newRegistrations, reg)
			}
		}

		tm.registrations[triggerType] = newRegistrations
	}
}

// ActivateTrigger activates all registered triggers of a specific type
func (tm *TriggerManager) ActivateTrigger(triggerType TriggerType, ctx TriggerContext) {
	registrations := tm.registrations[triggerType]

	// Create a temporary copy to avoid issues with modifications during iteration
	triggersToActivate := make([]TriggerRegistration, len(registrations))
	copy(triggersToActivate, registrations)

	logger.Debug("Activating triggers", logger.String("triggerType", triggerType.String()), logger.Int("count", len(triggersToActivate)))

	for _, reg := range triggersToActivate {
		reg.Callback(&ctx, reg.RegisteredOn)

		// If the trigger is one-time only, remove it after activation
		if reg.OneTimeOnly {
			tm.UnregisterTrigger(reg.ID)
		}
	}
}
