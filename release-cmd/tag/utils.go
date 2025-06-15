package tag

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type SetupConfig struct {
	BaseRelease   string `yaml:"base_tag"`
	FinalRelease  string `yaml:"target_tag"`
	OperationType string `yaml:"operation_type"`
}

var SetupConfigobject SetupConfig

func init() {
	errLoad := loadYaml("./release-cmd/promotional.yaml")
	log.Printf("Loading configuration from ./release-cmd/promotional.yaml")
	if errLoad != nil {
		log.Fatal(errLoad)
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

func GitPush(branch string) error {
	cmd := exec.Command("git", "push", "origin", branch)
	log.Printf("Running git push command: %s'", cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CleanWorkingDirectory() error {
	cmdReset := exec.Command("git", "reset", "--hard")
	cmdReset.Stdout = os.Stdout
	cmdReset.Stderr = os.Stderr
	if err := cmdReset.Run(); err != nil {
		return fmt.Errorf("failed to reset changes: %w", err)
	}
	cmdClean := exec.Command("git", "clean", "-fd")
	cmdClean.Stdout = os.Stdout
	cmdClean.Stderr = os.Stderr
	if err := cmdClean.Run(); err != nil {
		return fmt.Errorf("failed to clean untracked files: %w", err)
	}
	log.Println("Working directory cleaned successfully.")
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
		cmd := exec.Command("git", "tag", "-a", newTag, "-m", "Tagging Rc Branch")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create git tag: %v, output: %s", err, string(output))
		}
		log.Printf("Successfully created tag %s\n\n", newTag)
		return nil
	} else {
		cmd := exec.Command("git", "tag", newTag, existingTag)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create git tag: %v, output: %s", err, string(output))
		}
		log.Printf("Successfully created tag %s from %s\n\n", newTag, existingTag)
		return nil
	}
}

func FetchReleaseBranch(BranchName string) (string, error) {
	cmd := exec.Command("git", "branch", "-r", "--list", "origin/release/"+BranchName)
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

func ReleaseGithub(preRelease bool, tagName string) {
	repo := os.Getenv("GITHUB_REPOSITORY")
	token := os.Getenv("GH_PAT")
	cmdArgs := []string{"release", "create", "-R", fmt.Sprintf("https://github.com/%s", repo), "--generate-notes"}
	if preRelease == true {
		cmdArgs = append(cmdArgs, "--prerelease")
	}
	cmdArgs = append(cmdArgs, tagName)
	cmd := exec.Command("gh", cmdArgs...)
	cmd.Env = append(os.Environ(), fmt.Sprintf("GITHUB_TOKEN=%s", token))
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error: %v\n", err)
		log.Printf("Output: %s\n", string(output))
		return
	}
	log.Println("Release created successfully!")
}

func FetchTag(tagPattern string) string {
	cmd := exec.Command("git", "tag", "-l", "--sort=-v:refname", tagPattern)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error executing git command:", err)
		os.Exit(1)
	}
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(tags) == 0 || tags[0] == "" {
		log.Printf("No tags found matching the pattern: %s", tagPattern)
		return ""
	}
	return tags[0]
}

func NewSprintName() string {
	currentDate := time.Now()
	weekNumber := fmt.Sprintf("%d", (currentDate.Day()-1)/7+1)
	currYear := fmt.Sprintf("%02d", currentDate.Year()%100)
	currMonth := fmt.Sprintf("%02d", currentDate.Month())
	output := currYear + currMonth + weekNumber
	return output
}

func FetchTagToCheckIfItExists(tagName string) bool {
	cmd := exec.Command("git", "tag", "-l", tagName)
	output, err := cmd.Output()
	log.Printf("Output: %s\n", string(output))
	if err != nil {
		log.Printf("Error executing git command: %v\n", err)
		os.Exit(1)
	}

	if output == nil {
		log.Printf("No tags found matching the pattern: %s\n", tagName)
		return true
	} else {
		log.Printf("Tag %s found matching the pattern.\n", tagName)
		return false
	}
}

func CheckForHFfinalName(obj *SetupConfig) bool {

	if len(obj.BaseRelease) != 9 {
		log.Printf("HF base name %s is not valid, should be 9 characters long.\n", obj.BaseRelease)
		return false
	}

	num, err := strconv.Atoi(obj.BaseRelease)
	if err != nil {
		log.Printf("HF base name %s is not a valid number: %v\n", obj.BaseRelease, err)
		return false
	}

	finalNum, errFinal := strconv.Atoi(obj.FinalRelease)
	if errFinal != nil {
		log.Printf("HF final name %s is not a valid number: %v\n", obj.FinalRelease, errFinal)
		return false
	}

	if finalNum != num+1 {
		log.Printf("HF final name %s is not valid, should be one greater than base release %s.\n", obj.FinalRelease, obj.BaseRelease)
		return false
	} else {
		log.Printf("HF final name %s is valid, should be one greater than base release %s.\n", obj.FinalRelease, obj.BaseRelease)
		return true
	}

}
