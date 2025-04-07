package game

// TagType represents the different types of tags that can be applied to entities
type TagType int

// Tag constants representing various entity attributes and states
const (
	TAG_NONE TagType = iota
	TAG_TAUNT
	TAG_DIVINE_SHIELD
	TAG_CHARGE
	TAG_FROZEN
	TAG_STEALTH
	TAG_POISONOUS
	TAG_WINDFURY
	TAG_DEATHRATTLE
	TAG_BATTLECRY
	TAG_RUSH
	TAG_LIFESTEAL
	TAG_REBORN
	TAG_DORMANT
	TAG_SPELLPOWER
	TAG_CANT_ATTACK
	TAG_CANT_BE_TARGETED
	TAG_IMMUNE
)

// Tag represents a key-value pair for entity attributes in Hearthstone
type Tag struct {
	Type  TagType
	Value interface{} // Can hold different types (bool, int, etc.) depending on the tag
}

// NewTag creates a new tag with the specified type and value
func NewTag(tagType TagType, value interface{}) Tag {
	return Tag{
		Type:  tagType,
		Value: value,
	}
}

// Common tag helper functions
func HasTag(tags []Tag, tagType TagType) bool {
	for _, tag := range tags {
		if tag.Type == tagType {
			return true
		}
	}
	return false
}

// GetTagValue returns the value of a specific tag if it exists
func GetTagValue(tags []Tag, tagType TagType) (interface{}, bool) {
	for _, tag := range tags {
		if tag.Type == tagType {
			return tag.Value, true
		}
	}
	return nil, false
}

