package loader

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	Paths struct {
		ThemesDirectory     string `toml:"themes_directory"`
		AlacrittyConfigPath string `toml:"alacritty_config_path"`
	} `toml:"paths"`
	Repos struct {
		ThemeURL string `toml:"theme_url"`
	} `toml:"repos"`
}

// LoadConfig reads a TOML file and returns a Config instance.
func LoadConfig(path string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := toml.NewDecoder(file).Decode(config); err != nil {
		return nil, err
	}
	config.Paths.AlacrittyConfigPath = expandHome(config.Paths.AlacrittyConfigPath)
	config.Paths.ThemesDirectory = expandHome(config.Paths.ThemesDirectory)

	return config, nil
}

func expandHome(path string) string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Cannot get current user:", err)
		return path // return original path which might lead to further errors
	}
	if strings.HasPrefix(path, "~") {
		return filepath.Join(usr.HomeDir, path[2:])
	}
	if strings.HasPrefix(path, "$HOME") {
		return filepath.Join(usr.HomeDir, path[6:])
	}
	return path
}
