package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	cf "goalacritty_themes/config"
	it "goalacritty_themes/theme_tools"
	"os"
)

var (
	spinnerMod   = spinner.MiniDot
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

// spinnerModel and its Init, Update and View functions
// for cloning theme repo (if missing)
type spinnerModel struct {
	spinner spinner.Model
	config  cf.Config
	done    chan struct{}
}

// TODO: error handling
func (m spinnerModel) Init() tea.Cmd {
	go func() {
		err := it.InstallThemes(m.config)
		if err != nil {
			fmt.Println("Error installing themes:", err)
			os.Exit(1)
		}
		themes, err := it.GetThemeDataNames(m.config)
		if err != nil {
			fmt.Println("Error retrieving theme names:", err)
			os.Exit(1)
		}
		firstTheme := themes[0]
		err = it.InitAlacrittyConfig(m.config, firstTheme)
		if err != nil {
			fmt.Println("Error retrieving theme names:", err)
			os.Exit(1)
		}
		close(m.done)
	}()
	return m.spinner.Tick
}

func (m spinnerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case <-m.done:
		mainModel := InitializeMainModel(m.config)
		return mainModel, nil
	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}
	return m, nil
}

func (m spinnerModel) View() string {
	return fmt.Sprintf("\n\n   %s Installing themes...", m.spinner.View())
}

func InitializeSpinnerModel(config cf.Config) spinnerModel {
	s := spinner.New()
	s.Spinner = spinnerMod
	s.Style = spinnerStyle

	m := spinnerModel{
		spinner: s,
		config:  config,
		done:    make(chan struct{}),
	}
	return m
}
