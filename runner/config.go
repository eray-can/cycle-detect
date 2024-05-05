package runner

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path/filepath"
)

const (
	DetectToml = "detect.toml"
)

type Config struct {
	Build *Build `toml:"build"`
}

type Log struct {
	ShowScreen bool `toml:"show_screen"`
}

type Build struct {
	ExcludeDir      []*string `toml:"exclude_dir"`
	BaseProjectPath *string   `toml:"project_path"`
}

func ReadFileConfig() (*Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, errors.New("detect.toml not found read to default conf")
	}

	file, fileErr := os.ReadFile(filepath.Join(*configPath, DetectToml))
	if fileErr != nil {
		return nil, fmt.Errorf("%s detect.toml file not read", DetectToml)
	}

	var config Config
	if err = toml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("detect.toml not unmarshalled %s", err.Error())
	}

	checkErr := configChecker(&config)
	if checkErr != nil {
		return nil, checkErr
	}

	return &config, nil
}

func getConfigFilePath() (*string, error) {
	basePath, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}
	return &basePath, nil
}

func configChecker(config *Config) error {
	if config.Build.BaseProjectPath == nil || *config.Build.BaseProjectPath == "" {
		return fmt.Errorf("%s [build][baseProjectPath] not found", DetectToml)
	}

	return nil
}
