package game

import (
	"testing"
)

// TestNewPlayerManaInitialization tests that NewPlayer correctly initializes mana values
func TestNewPlayerManaInitialization(t *testing.T) {
	player := NewPlayer()

	// Verify MaxMana is initialized correctly
	if player.MaxMana != DefaultMaxMana {
		t.Errorf("Expected MaxMana to be %d, got %d", DefaultMaxMana, player.MaxMana)
	}

	// Verify Mana is initialized correctly
	if player.Mana != DefaultStartingMana {
		t.Errorf("Expected Mana to be %d, got %d", DefaultStartingMana, player.Mana)
	}

	// Verify TotalMana is initialized correctly
	if player.TotalMana != DefaultStartingMana {
		t.Errorf("Expected TotalMana to be %d, got %d", DefaultStartingMana, player.TotalMana)
	}
}
