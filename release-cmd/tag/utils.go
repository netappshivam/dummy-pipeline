package tag

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
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
	log.Printf("Running 'git checkout %s'", cmd)
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
	log.Printf("Running %s'", cmd)
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
	if existingTag == "" {
		cmd := exec.Command("git", "tag", "-a", newTag, "-m", "Tagging feature branch")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create git tag: %v, output: %s", err, string(output))
		}
		fmt.Printf("Successfully created tag %s\n", newTag)
		return nil
	} else {
		cmd := exec.Command("git", "tag", newTag, existingTag)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create git tag: %v, output: %s", err, string(output))
		}
		fmt.Printf("Successfully created tag %s from %s\n", newTag, existingTag)
		return nil
	}
}

func ReleaseGithub() {

	// Construct the command
	repo := os.Getenv("GITHUB_REPOSITORY")
	token := os.Getenv("GH_PAT")

	cmdArgs := []string{"release", "create", "-R", fmt.Sprintf("https://github.com/%s", repo), "--generate-notes"}
	if strings.Contains(SetupConfigobject.BaseRelease, "-DEV.") {
		cmdArgs = append(cmdArgs, "--prerelease")
	}
	cmdArgs = append(cmdArgs, SetupConfigobject.FinalRelease)

	// Execute the command
	cmd := exec.Command("gh", cmdArgs...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", token))

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return
	}

	fmt.Println("Release created successfully!")
}

func FetchDevTag() string {
	tagPattern := "*-DEV.*"

	cmd := exec.Command("git", "tag", "-l", "--sort=-v:refname", tagPattern)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing git command:", err)
		os.Exit(1)
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	log.Printf(string(output))

	if len(tags) == 0 || tags[0] == "" {
		log.Printf("No tags found matching the pattern:", tagPattern)
		os.Exit(0)
	}

	return tags[0]

}

func NewSprintName() string {

	currentDate := time.Now()
	// Calculate the week number of the month
	weekNumber := fmt.Sprintf("%d", (currentDate.Day()-1)/7+1)
	currYear := fmt.Sprintf("%02d", currentDate.Year()%100)
	currMonth := fmt.Sprintf("%02d", currentDate.Month())

	output := currYear + currMonth + weekNumber

	return output

}
