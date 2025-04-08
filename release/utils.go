package release_branch

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
)

type SetupConfig struct {
	current_week_release string `yaml:"current_week_release"`
	next_week_release    string `yaml:"next_week_release"`
}

var setupConfig SetupConfig

func intializeData() {

	loadYaml("release.yaml")

}

func loadYaml(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening .yaml file: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading .yaml file: %v", err)
	}

	err = yaml.Unmarshal(data, &setupConfig)
	if err != nil {
		return fmt.Errorf("error unmarshalling .yaml file: %v", err)
	}
	return nil
}
