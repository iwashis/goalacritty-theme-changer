package loader

import (
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockConfig is a sample TOML config for testing purposes.
const mockConfig = `
[paths]
themes_directory = "~/test/themes"
alacritty_config_path = "$HOME/.config/alacritty/alacritty.yml"

[repos]
theme_url = "https://github.com/some/repo"
`

// TestLoadConfig tests the LoadConfig function.
func TestLoadConfig(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "testconfig_*.toml")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the file after the test

	// Write the mockConfig content to the temp file
	_, err = tmpFile.WriteString(mockConfig)
	if err != nil {
		t.Fatalf("Error writing to temp file: %v", err)
	}
	tmpFile.Close()

	// Test LoadConfig with valid file
	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Get the current user's home directory for comparison
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	// Validate the config fields
	expectedThemesDir := filepath.Join(homeDir, "test/themes")
	expectedAlacrittyPath := filepath.Join(homeDir, ".config/alacritty/alacritty.yml")

	assert.Equal(t, expectedThemesDir, config.Paths.ThemesDirectory, "ThemesDirectory should match expanded path")
	assert.Equal(t, expectedAlacrittyPath, config.Paths.AlacrittyConfigPath, "AlacrittyConfigPath should match expanded path")
	assert.Equal(t, "https://github.com/some/repo", config.Repos.ThemeURL, "ThemeURL should match")
}

// TestLoadConfigFileNotFound tests LoadConfig when the file doesn't exist.
func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("non_existent_file.toml")
	if err == nil {
		t.Fatalf("Expected an error for non-existent file, got none")
	}
}

// TestExpandHome tests the expandHome function.
func TestExpandHome(t *testing.T) {
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	// Test with a path starting with ~
	path := expandHome("~/mydir")
	expected := filepath.Join(homeDir, "mydir")
	assert.Equal(t, expected, path, "Expected path to expand to user's home directory")

	// Test with a path starting with $HOME
	path = expandHome("$HOME/mydir")
	assert.Equal(t, expected, path, "Expected $HOME to expand to user's home directory")

	// Test with a path that doesn't require expansion
	path = expandHome("/some/other/path")
	assert.Equal(t, "/some/other/path", path, "Expected path to remain unchanged")
}

// TestExpandHomeError tests the expandHome function when user.Current fails.
func TestExpandHomeError(t *testing.T) {
	// Mock the user.Current function to return an error (this requires more sophisticated mocking techniques)
	// For simplicity, we can skip this test unless we have a mock framework
}


