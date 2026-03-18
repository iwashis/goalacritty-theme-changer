package main

import models "goalacritty_themes/models"
import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	cf "goalacritty_themes/config"
	it "goalacritty_themes/theme_tools"
	"os"
  "path/filepath"
)

func main() {
  // Step 1. Define paths to check
	localConfigPath := "config.toml"
	homeConfigPath := filepath.Join(os.Getenv("HOME"), ".config", "goalacritty-themes", "config.toml")

	// Step 2. Attempt to load config from local directory first
	config, err := cf.LoadConfig(localConfigPath)
	if err != nil {
		fmt.Println("Local config not found, trying global config:", err)
		
		// Step 3. Fallback to global config in $HOME/.config/goalacritty-themes/
		config, err = cf.LoadConfig(homeConfigPath)
		if err != nil {
			fmt.Println("Error loading config from global path:", err)
			return
		}
	}
	// Check if theme repo is in place
	if !it.IsThemesRepoInstalled(*config) {
		// if the repo is missing install it using spinnerModel bubbletea functionality
		m := models.InitializeSpinnerModel(*config)
		if _, err := tea.NewProgram(m).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
		return
	} else {
    // here, the repository is in place. Run the main model
    mainModel := models.InitializeMainModel(*config)
    if _, err := tea.NewProgram(mainModel).Run(); err != nil {
      fmt.Println("Error running program:", err)
      os.Exit(1)
    }
    return
  }
}
