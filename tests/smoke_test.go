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

func TestEngineAutoRun(t *testing.T) {
	// Initialize the application with test config
	configPath := filepath.Join("testdata", "smoke_test.json")
	if err := bootstrap.Initialize(configPath); err != nil {
		t.Fatalf("Failed to initialize global components: %v", err)
	}

	// Load a game
	gameManager := game.GetGameManager()
	g, err := gameManager.LoadGameByID("smoke_test")
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Create a game engine with auto-run enabled (default)
	e := engine.NewEngine(g)

	// Start the game - it should automatically progress to the MainAction phase
	err = e.StartGame()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Verify that auto-run progressed to MainAction phase
	if g.Phase != game.MainAction {
		t.Errorf("Expected auto-run to progress to MainAction phase, but got %v", g.Phase)
	}

	// Verify that the game turn is 1
	if g.CurrentTurn != 1 {
		t.Errorf("Expected turn to be 1 after game start, got %d", g.CurrentTurn)
	}

	// Test end turn
	err = e.EndPlayerTurn()
	if err != nil {
		t.Fatalf("Failed to end turn: %v", err)
	}

	// Verify phase progression after end turn (should be in MainAction of next turn)
	if g.Phase != game.MainAction {
		t.Errorf("Expected to be in MainAction phase after ending turn, but got %v", g.Phase)
	}

	// Verify turn has been incremented
	if g.CurrentTurn != 2 {
		t.Errorf("Expected turn to be 2 after ending turn, got %d", g.CurrentTurn)
	}

	// Test with auto-run disabled
	// Reset the game
	g, err = gameManager.LoadGameByID("smoke_test")
	if err != nil {
		t.Fatalf("Failed to load game: %v", err)
	}

	// Create a new engine
	e = engine.NewEngine(g)
	
	// Disable auto-run
	e.SetAutoRun(false)

	// Start game - it should stay at InvalidPhase with engine.nextPhase set to BeginFirst
	err = e.StartGame()
	if err != nil {
		t.Fatalf("Failed to start game: %v", err)
	}

	// Test manual progression
	// Process next phase - should move to BeginFirst
	err = e.ProcessNextPhase()
	if err != nil {
		t.Fatalf("Failed to process next phase: %v", err)
	}

	// Verify that game stopped at BeginFirst
	if g.Phase != game.BeginFirst {
		t.Errorf("Expected game to stop at BeginFirst with auto-run disabled, but got %v", g.Phase)
	}

	// Test ProcessUntil
	err = e.ProcessUntil(game.MainAction)
	if err != nil {
		t.Fatalf("Failed to process until MainAction: %v", err)
	}

	// Verify phase is now MainAction
	if g.Phase != game.MainAction {
		t.Errorf("Expected phase to be MainAction after ProcessUntil, but got %v", g.Phase)
	}
}
