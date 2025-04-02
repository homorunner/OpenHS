package game

import (
	"testing"

	"github.com/openhs/internal/config"
)

// TestNewPlayerManaInitialization tests that NewPlayer correctly initializes mana values
func TestNewPlayerManaInitialization(t *testing.T) {
	player := NewPlayer()

	// Verify MaxMana is initialized correctly
	if player.MaxMana != config.DefaultMaxMana {
		t.Errorf("Expected MaxMana to be %d, got %d", config.DefaultMaxMana, player.MaxMana)
	}

	// Verify Mana is initialized correctly
	if player.Mana != config.DefaultStartingMana {
		t.Errorf("Expected Mana to be %d, got %d", config.DefaultStartingMana, player.Mana)
	}

	// Verify TotalMana is initialized correctly
	if player.TotalMana != config.DefaultStartingMana {
		t.Errorf("Expected TotalMana to be %d, got %d", config.DefaultStartingMana, player.TotalMana)
	}
} 