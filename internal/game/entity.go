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

// Zone represents different zones in the game
type Zone int

const (
	ZONE_NONE Zone = iota
	ZONE_PLAY
	ZONE_DECK
	ZONE_HAND
	ZONE_GRAVEYARD
	ZONE_REMOVEDFROMGAME
	ZONE_SETASIDE
	ZONE_SECRET
)

// String returns a string representation of the Zone
func (z Zone) String() string {
	switch z {
	case ZONE_PLAY:
		return "Play"
	case ZONE_DECK:
		return "Deck"
	case ZONE_HAND:
		return "Hand"
	case ZONE_GRAVEYARD:
		return "Graveyard"
	case ZONE_REMOVEDFROMGAME:
		return "RemovedFromGame"
	case ZONE_SETASIDE:
		return "SetAside"
	case ZONE_SECRET:
		return "Secret"
	default:
		return "None"
	}
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
