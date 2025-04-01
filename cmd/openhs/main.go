package main

import (
	"fmt"
	"path/filepath"

	"github.com/openhs/internal/bootstrap"
)

func main() {
	configPath := filepath.Join("config", "openhs.json")
	if err := bootstrap.Initialize(configPath); err != nil {
		fmt.Printf("Failed to initialize global components: %v\n", err)
		return
	}

	fmt.Print("OpenHS - Hearthstone Simulator Core")
}