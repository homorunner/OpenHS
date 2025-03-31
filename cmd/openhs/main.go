package main

import (
	"fmt"

	"github.com/openhs/internal/logger"
)

func main() {
	// Initialize logger with debug level for development
	_, err := logger.SetLogger("debug", "logs", "openhs.log", true)
	if err != nil {
		logger.Fatal("Failed to initialize logger", logger.Err(err))
	}

	fmt.Print("OpenHS - Hearthstone Simulator Core")
}