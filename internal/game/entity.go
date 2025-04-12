package game

// Entity represents a card instance in play with a reference to its definition and owner
type Entity struct {
	Card              *Card
	Owner             *Player
	Health            int
	MaxHealth         int
	Attack            int
	Tags              []Tag  // Store entity states like Taunt, Divine Shield, etc.
	Buffs             []Buff // Track any modifications specific to this instance
	IsDestroyed       bool
	NumAttackThisTurn int  // Tracks how many times this entity has attacked this turn
	Exhausted         bool // Indicates if the entity can attack or not
	NumTurnInPlay     int  // Tracks how many turns the entity has been in field (0 = first turn)
	CurrentZone       Zone // Tracks which zone the entity is in
}

// NewEntity creates a new entity from a card
func NewEntity(card *Card, game *Game, owner *Player) *Entity {
	entity := &Entity{
		Card:        card,
		Owner:       owner,
		Health:      card.Health,
		MaxHealth:   card.Health,
		Attack:      card.Attack,
		Tags:        make([]Tag, 0, len(card.Tags)), // Preallocate capacity
		Buffs:       make([]Buff, 0),
		CurrentZone: ZONE_NONE, // Initial zone is NONE until placed somewhere
	}

	// Copy tags from card to entity
	entity.Tags = append(entity.Tags, card.Tags...)

	// Load card effects
	if card.Load != nil {
		card.Load(game, entity)
	}

	return entity
}
