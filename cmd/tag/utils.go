package tag

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type SetupConfig struct {
	CurrentWeekRelease string `yaml:"current_weekly_release"`
	NextWeekRelease    string `yaml:"next_weekly_release"`
}

var SetupConfigobject SetupConfig

func init() {

	loadYaml("./cmd/releases.yaml")
}

func IncrementTag(tag string) (string, error) {
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

func FetchTags() error {
	cmd := exec.Command("git", "fetch", "--tags")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch tags: %s, output: %s", err, string(output))
	}
	return nil
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

func FetchReleaseBranch(sprint string) (string, error) {
	cmd := exec.Command("git", "branch", "-r", "--list", "origin/release."+sprint)
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
	defer file.Close()
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
func GithubUserEmail() {
	user, err := github.GetGithubUser(github.GhToken, github.PrUser)
	if err != nil {
		log.Println("Error:", err)
		os.Exit(1)
	}

	if user.Email == nil {
		log.Println("Error: Email not available for user:", github.PrUser)
		os.Exit(1)
	}

	ghEmail := *user.Email
}
