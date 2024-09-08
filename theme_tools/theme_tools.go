package install_themes

import configloader "goalacritty_themes/config"
import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
)

// Function to check if the themes repository is already installed
func IsThemesRepoInstalled(config configloader.Config) bool {
	// Get the path to the themes directory
	themesDir := os.ExpandEnv(config.Paths.ThemesDirectory)

	// Check if the themes directory exists
	if _, err := os.Stat(themesDir); os.IsNotExist(err) {
		return false
	}

	// Check if the directory is a Git repository by looking for the .git directory
	gitDir := filepath.Join(themesDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return false
	}
	return true
}

func InstallThemes(config configloader.Config) error {
	// Step 1: Create the directory
	createDirCmd := exec.Command("mkdir", "-p", os.ExpandEnv(config.Paths.ThemesDirectory))
	createDirCmd.Env = os.Environ()
	if err := createDirCmd.Run(); err != nil {
    return err
	}

	// Step 2: Clone the GitHub repository
	cloneRepoCmd := exec.Command("git", "clone", config.Repos.ThemeURL, os.ExpandEnv(config.Paths.ThemesDirectory))
	cloneRepoCmd.Env = os.Environ()
	if err := cloneRepoCmd.Run(); err != nil {
		return err
	}
	fmt.Println("Alacritty-theme repository cloned successfully")
  return nil
}

func handleExecError(err error) {
	if exitError, ok := err.(*exec.ExitError); ok {
		if status, ok := exitError.Sys().(syscall.WaitStatus); ok && status.ExitStatus() == 128 {
			fmt.Println("It appears that you already have alacritty-theme installed on your system")
			return
		}
	}
	fmt.Println("Error encountered:", err)
}

// ThemeData represents a theme file with its name and full path
type ThemeData struct {
	Name     string
	FullPath string
}

// GetThemeDataNames reads the names and full paths of files in the config.Paths.Themes directory
func GetThemeDataNames(config configloader.Config) ([]ThemeData, error) {
	dir := os.ExpandEnv(config.Paths.ThemesDirectory + "/themes/")
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	var themeFiles []ThemeData
	for _, file := range files {
		if !file.IsDir() {
			fullPath := filepath.Join(dir, file.Name())
			themeFiles = append(themeFiles, ThemeData{
				Name:     strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
				FullPath: fullPath,
			})
		}
	}
	return themeFiles, nil
}

func GetCurrentTheme(config configloader.Config) (*ThemeData, error) {

	alacrittyConfigPath := config.Paths.AlacrittyConfigPath
	expandedThemesDirectory := config.Paths.ThemesDirectory
	themePathPattern := regexp.QuoteMeta(expandedThemesDirectory+"/themes/") + `(\w+)\.toml`

	// Read the Alacritty config file
	content, err := os.ReadFile(alacrittyConfigPath) // Use expanded path
	if err != nil {
		return nil, err
	}

	// Replace the old theme path with the new one using a regex
	re := regexp.MustCompile(themePathPattern)
	currentThemePath := re.FindString(string(content))
	return &ThemeData{
		Name:     currentThemePath,
		FullPath: currentThemePath,
	}, nil
}


// InitAlacrittyConfig appends the themePath at the bottom of the configuration file if it is not already present.
func InitAlacrittyConfig(config configloader.Config, theme ThemeData) error {
	// Read the existing Alacritty config file
  alacrittyConfigPath := config.Paths.AlacrittyConfigPath
  themePath := theme.FullPath
	themePattern := regexp.QuoteMeta("/themes/") + `(\w+)\.toml`

 // Check if the Alacritty config file exists
  if _, err := os.Stat(alacrittyConfigPath); os.IsNotExist(err) {
      // Create the file if it does not exist
      file, err := os.Create(alacrittyConfigPath)
      if err != nil {
          return fmt.Errorf("failed to create the configuration file: %w", err)
      }
      file.Close() // Close the file after creating it
  } else if err != nil {
      return fmt.Errorf("failed to check if the configuration file exists: %w", err)
  }

  // Reading the Alacritty config file:
	content, err := os.ReadFile(alacrittyConfigPath)
	if err != nil {
		return err // Return any error encountered while reading
	}

	// Convert the content to a string for regex matching
	contentStr := string(content)

	// Check if the theme path already exists in the file
	re := regexp.MustCompile(themePattern)

	// If the theme path is found, do nothing and exit
	if re.MatchString(contentStr) {
		return nil
  } 

	// If the theme path is not found, append it to the content
	contentStr = "import = [\n" + "\"" + themePath + "\"\n" + "]\n\n" + contentStr

	// Write the modified content back to the config file
	return os.WriteFile(alacrittyConfigPath, []byte(contentStr), 0644)
}


func UpdateAlacrittyConfigFile(config configloader.Config, td ThemeData) error {
	alacrittyConfigPath := config.Paths.AlacrittyConfigPath
	expandedThemesDirectory := config.Paths.ThemesDirectory
	newThemePath := td.FullPath

	// Construct the old and new theme paths to search and replace
	themePattern := regexp.QuoteMeta(expandedThemesDirectory + "/themes/") + `[^\"]+\.toml`

	// Read the Alacritty config file
	content, err := os.ReadFile(alacrittyConfigPath)
	if err != nil {
		return err
	}

	// Check if the old theme path exists and replace it
	re := regexp.MustCompile(themePattern)
	if re.MatchString(string(content)) {
		// Replace the old theme path with the new one
		modifiedContent := re.ReplaceAllString(string(content), newThemePath)
		return os.WriteFile(alacrittyConfigPath, []byte(modifiedContent), 0644)
	}

	// If no match is found, append the new theme path
	return InitAlacrittyConfig(config, td)
}

