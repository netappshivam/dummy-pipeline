package tag

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type SetupConfig struct {
	BaseRelease  string `yaml:"base_tag"`   //"2512.0.0-DEV.0"
	FinalRelease string `yaml:"target_tag"` //"2512.0.0"
}

var SetupConfigobject SetupConfig

func init() {

	err_load := loadYaml("./release-cmd/releases.yaml")
	if err_load != nil {
		log.Fatal(err_load)
		return
	}

}

func FetchTagsPrune() error {
	cmd := exec.Command("git", "fetch", "--prune", "--tags")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %s, output: %s", err, string(output))
	}
	return nil
}

func GitCheckout(branch, ref string) error {
	cmd := exec.Command("git", "checkout", "-b", branch, ref)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GitOnlyCheckout(branch string) error {
	cmd := exec.Command("git", "checkout", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GitPush(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func FetchReleaseBranch(BranchName string) (string, error) {
	cmd := exec.Command("git", "branch", "-r", "--list", "origin/"+BranchName)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return "", nil
	}
	return branch, nil
}

func CleanWorkingDirectory() error {
	// Discard local changes
	cmdReset := exec.Command("git", "reset", "--hard")
	cmdReset.Stdout = os.Stdout
	cmdReset.Stderr = os.Stderr
	if err := cmdReset.Run(); err != nil {
		return fmt.Errorf("failed to reset changes: %w", err)
	}

	// Remove untracked files and directories
	cmdClean := exec.Command("git", "clean", "-fd")
	cmdClean.Stdout = os.Stdout
	cmdClean.Stderr = os.Stderr
	if err := cmdClean.Run(); err != nil {
		return fmt.Errorf("failed to clean untracked files: %w", err)
	}

	fmt.Println("Working directory cleaned successfully.")
	return nil
}

func loadYaml(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("error opening .yaml file: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading .yaml file: %v", err)
	}

	err = yaml.Unmarshal(data, &SetupConfigobject)
	if err != nil {
		return fmt.Errorf("error unmarshalling .yaml file: %v", err)
	}
	return nil
}

func CreateGitTag(newTag, existingTag string) error {
	cmd := exec.Command("git", "tag", newTag, existingTag)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create git tag: %v, output: %s", err, string(output))
	}
	fmt.Printf("Successfully created tag %s from %s\n", newTag, existingTag)
	return nil
}
