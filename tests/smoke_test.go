package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/openhs/internal/bootstrap"
	"github.com/openhs/internal/game"
)

// TestGameLoadingAndInitialization tests that games can be loaded from configuration
func TestGameLoadingAndInitialization(t *testing.T) {
	// Run in the root directory
	os.Chdir("..")

	// Initialize the application with test config
	configPath := filepath.Join("config", "openhs.json")
	if err := bootstrap.Initialize(configPath); err != nil {
		t.Fatalf("Failed to initialize global components: %v", err)
	}

	// Load a game directly using the game manager
	gameManager := game.GetGameManager()
	g, err := gameManager.LoadGameByID("smoke_test")
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Verify game was loaded correctly
	if g == nil {
		t.Fatal("Game is nil after loading")
	}
	if len(g.Players) != 2 {
		t.Errorf("Expected 2 players, got %d", len(g.Players))
	}

	// Verify game state after initialization
	if g.CurrentTurn != 0 {
		t.Errorf("Expected initial turn to be 0, got %d", g.CurrentTurn)
	}
	if g.Phase != game.InvalidPhase {
		t.Errorf("Expected initial phase to be InvalidPhase, got %v", g.Phase)
	}
}
