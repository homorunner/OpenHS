package tests

import (
	"path/filepath"
	"testing"

	"github.com/openhs/internal/bootstrap"
	"github.com/openhs/internal/engine"
	"github.com/openhs/internal/game"
)

func TestGameLoadingAndInitialization(t *testing.T) {
	// Initialize the application with test config
	configPath := filepath.Join("testdata", "smoke_test.json")
	if err := bootstrap.Initialize(configPath); err != nil {
		t.Fatalf("Failed to initialize global components: %v", err)
	}

	// Load a game directly using the game manager
	gameManager := game.GetGameManager()
	g, err := gameManager.LoadGameByID("sample_game")
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

	// Create a game engine
	e := engine.NewEngine(g)
	if e == nil {
		t.Fatal("Failed to create game engine")
	}

	// Verify game state after initialization
	if g.CurrentTurn != 0 {
		t.Errorf("Expected initial turn to be 0, got %d", g.CurrentTurn)
	}
	if g.Phase != game.StartGame {
		t.Errorf("Expected initial phase to be StartGame, got %v", g.Phase)
	}
}
