package main

import models "goalacritty_themes/models"
import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	cf "goalacritty_themes/config"
	it "goalacritty_themes/theme_tools"
	"os"
)

func main() {
	// Step 1. load config
	config, err := cf.LoadConfig("config.toml")
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
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
