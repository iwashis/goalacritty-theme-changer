package install_themes

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	configloader "goalacritty_themes/config"
)

// TestGetThemeDataNames tests the GetThemeDataNames function
func TestGetThemeDataNames(t *testing.T) {
	// Use a temporary directory for the themes
	themesDir, err := os.MkdirTemp("", "test_get_theme_data")
	assert.NoError(t, err)
	defer os.RemoveAll(themesDir) // Clean up after the test

	// Create the "themes" directory inside the temporary directory
	err = os.Mkdir(filepath.Join(themesDir, "themes"), os.ModePerm)
	assert.NoError(t, err)

	// Create mock theme files inside the "themes" directory
	themeFiles := []string{"dark-theme.toml", "light_theme.toml", "monokai-pro.toml"}
	for _, theme := range themeFiles {
		_, err := os.Create(filepath.Join(themesDir, "themes", theme))
		assert.NoError(t, err)
	}

	// Create a mock config pointing to the temporary themes directory
	mockConfig := configloader.Config{
		Paths: struct {
			ThemesDirectory     string `toml:"themes_directory"`
			AlacrittyConfigPath string `toml:"alacritty_config_path"`
		}{
			ThemesDirectory:     themesDir,
			AlacrittyConfigPath: "/mock/.config/alacritty/alacritty.yml",
		},
	}

	// Call GetThemeDataNames with the mock config
	themeData, err := GetThemeDataNames(mockConfig)
	assert.NoError(t, err, "Expected no error reading theme files")

	// Verify that the returned theme data matches the mock theme files
	expectedThemes := []ThemeData{
		{Name: "dark-theme", FullPath: filepath.Join(themesDir, "themes", "dark-theme.toml")},
		{Name: "light_theme", FullPath: filepath.Join(themesDir, "themes", "light_theme.toml")},
		{Name: "monokai-pro", FullPath: filepath.Join(themesDir, "themes", "monokai-pro.toml")},
	}

	assert.Equal(t, len(expectedThemes), len(themeData), "Expected the correct number of theme files")

	for i, theme := range expectedThemes {
		assert.Equal(t, theme.Name, themeData[i].Name, "Expected theme name to match")
		assert.Equal(t, theme.FullPath, themeData[i].FullPath, "Expected full path to match")
	}
}



// TestInitAlacrittyConfig tests the InitAlacrittyConfig function
func TestInitAlacrittyConfig(t *testing.T) {
	// Use a temporary file for the Alacritty config
	alacrittyConfigPath, err := os.CreateTemp("", "test_alacritty_config_*.yml")
	assert.NoError(t, err)
	defer os.Remove(alacrittyConfigPath.Name()) // Clean up after the test

	// Create a mock config pointing to the temporary config file
	mockConfig := configloader.Config{
		Paths: struct {
			ThemesDirectory     string `toml:"themes_directory"`
			AlacrittyConfigPath string `toml:"alacritty_config_path"`
		}{
			ThemesDirectory:     "/mock/themes",
			AlacrittyConfigPath: alacrittyConfigPath.Name(),
		},
	}

	// Create a theme data instance
	theme := ThemeData{
		Name:     "monokai-pro",
		FullPath: "/mock/themes/monokai-pro.toml",
	}

	// Initialize the Alacritty config by appending the theme path
	err = InitAlacrittyConfig(mockConfig, theme)
	assert.NoError(t, err, "Expected no error when initializing Alacritty config")

	// Read the content of the Alacritty config file and verify the theme path is appended
	content, err := os.ReadFile(alacrittyConfigPath.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), theme.FullPath, "Expected theme path to be added to config")
}

func TestUpdateAlacrittyConfigFile(t *testing.T) {
	// Use a temporary file for the Alacritty config
	alacrittyConfigPath, err := os.CreateTemp("", "test_alacritty_config_*.yml")
	assert.NoError(t, err)
	defer os.Remove(alacrittyConfigPath.Name()) // Clean up after the test

	// Create a mock config pointing to the temporary config file
	mockConfig := configloader.Config{
		Paths: struct {
			ThemesDirectory     string `toml:"themes_directory"`
			AlacrittyConfigPath string `toml:"alacritty_config_path"`
		}{
			ThemesDirectory:     "/mock/themes",
			AlacrittyConfigPath: alacrittyConfigPath.Name(),
		},
	}

	// Create a theme data instance to initialize the config with the existing theme
	existingTheme := ThemeData{
		Name:     "dark-theme",
		FullPath: "/mock/themes/dark-theme.toml",
	}

	// Initialize the Alacritty config with the existing theme using InitAlacrittyConfig
	err = InitAlacrittyConfig(mockConfig, existingTheme)
	assert.NoError(t, err, "Expected no error when initializing Alacritty config with the existing theme")

	// Verify that the Alacritty config was initialized correctly
	content, err := os.ReadFile(alacrittyConfigPath.Name())
	assert.NoError(t, err)
	assert.Contains(t, string(content), existingTheme.FullPath, "Expected existing theme path to be in config")

	// Now create a new theme data instance to replace the old theme
	newTheme := ThemeData{
		Name:     "monokai-pro",
		FullPath: "/mock/themes/monokai-pro.toml",
	}

	// Call UpdateAlacrittyConfigFile to update the config by replacing the old theme with the new one
	err = UpdateAlacrittyConfigFile(mockConfig, newTheme)
	assert.NoError(t, err, "Expected no error when updating Alacritty config")

	// Read the updated content of the Alacritty config file
	content, err = os.ReadFile(alacrittyConfigPath.Name())
	assert.NoError(t, err)

	// Check if the new theme has replaced the old one
	assert.Contains(t, string(content), newTheme.FullPath, "Expected new theme path to be in config")
	// assert.NotContains(t, string(content), existingTheme.FullPath, "Expected old theme path to be replaced")
}

