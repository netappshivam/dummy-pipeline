package release_branch_logic

import (
	"fmt"
	//"gopkg.in/yaml.v2"
	//"io"
	//"os"
	"os/exec"
	"strconv"
	"strings"
)

type SetupConfig struct {
	current_week_release string
	next_week_release    string
}

var setupConfig SetupConfig

func IntializeData() {

	//loadYaml("release.yaml")
	setupConfig.current_week_release = "2515"
	setupConfig.next_week_release = "2516"

}

func incrementTag(tag string) (string, error) {
	// Split the tag into parts by "."
	parts := strings.Split(tag, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid tag format: %s", tag)
	}

	// Parse the last part as an integer
	lastPart := parts[3]
	lastNum, err := strconv.Atoi(lastPart)
	if err != nil {
		return "", fmt.Errorf("failed to parse last part of tag as number: %s", lastPart)
	}

	// Increment the last number
	lastNum++

	// Reconstruct the tag with the incremented number
	parts[3] = strconv.Itoa(lastNum)
	return strings.Join(parts, "."), nil
}

func fetchTags() error {
	cmd := exec.Command("git", "fetch", "--tags")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %s, output: %s", err, string(output))
	}
	return nil
}

func SetNewTag(newTag string) {
	fmt.Printf("::set-output name=new_tag::%s\n", newTag)
}

//func loadYaml(filepath string) error {
//	file, err := os.Open(filepath)
//	if err != nil {
//		return fmt.Errorf("error opening .yaml file: %v", err)
//	}
//	defer file.Close()
//	data, err := io.ReadAll(file)
//	if err != nil {
//		return fmt.Errorf("error reading .yaml file: %v", err)
//	}
//
//	err = yaml.Unmarshal(data, &setupConfig)
//	if err != nil {
//		return fmt.Errorf("error unmarshalling .yaml file: %v", err)
//	}
//	return nil
//}
