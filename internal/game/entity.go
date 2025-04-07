package game

// Entity represents a card instance in play with a reference to its definition and owner
type Entity struct {
	Card        *Card
	Owner       *Player
	Health      int
	MaxHealth   int
	Attack      int
	Tags        []Tag  // Store entity states like Taunt, Divine Shield, etc.
	Buffs       []Buff // Track any modifications specific to this instance
	IsDestroyed bool
}

// NewEntity creates a new entity from a card
func NewEntity(card *Card, owner *Player) *Entity {
	entity := &Entity{
		Card:      card,
		Owner:     owner,
		Health:    card.Health,
		MaxHealth: card.MaxHealth,
		Attack:    card.Attack,
		Tags:      make([]Tag, 0, len(card.Tags)), // Preallocate capacity
		Buffs:     make([]Buff, 0),
	}
	
	// Copy tags from card to entity
	entity.Tags = append(entity.Tags, card.Tags...)
	
	return entity
}