package runner

import "testing"

func TestConfigFilePath(t *testing.T) {
	_, err := getConfigFilePath()
	if err != nil {
		t.Error(err)
	}
}

func TestConfigChecker(t *testing.T) {
	err := configChecker(mockConfig())
	if err != nil {
		t.Error(err)
	}

}

func mockConfig() *Config {
	projectPath := dir()
	vendor := "vendor"

	return &Config{
		Build: &Build{
			BaseProjectPath: &projectPath,
			ExcludeDir: []*string{
				&vendor,
			},
		},
	}
}
