package game

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
