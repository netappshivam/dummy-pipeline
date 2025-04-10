package release_branch_logic

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Release_creation() {
	sprint := setupConfig.current_week_release
	fmt.Printf("Creating %s Branch\n", sprint)

	// finding if a branch exists for a weekly release
	last_branch, err := FetchReleaseBranch(sprint)
	if err != nil {
		log.Fatalf("Error fetching release branch: %v", err)
	}

	//if it does exists, then finding the latest branch tag and incrementing it
	if last_branch != "" {
		lastTag := LatestTagFetch("RC")
		fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

		newTag, _ := incrementTag(lastTag)

		SetNewTag(newTag)

		//logic here

	} else {
		//if it does not exist, then creating a new branch for that weekly release
		err := cleanWorkingDirectory()
		if err != nil {
			log.Fatalf("Error cleaning working directory: %v", err)
		}

		latestSprintTag := "release." + sprint

		fmt.Printf("Since there is no tag created for the sprint - %s, cutting the release main based out of master main\n", sprint)
		fmt.Printf("Creating the main - %s\n", sprint)
		if err := gitCheckout(latestSprintTag, "origin/main"); err != nil {
			log.Fatalf("Error creating main: %v", err)
		}
		if err := gitPush(latestSprintTag); err != nil {
			log.Fatalf("Error pushing main: %v", err)
		}
		SetNewTag(sprint + ".0.0-RC.1")
	}
}

func gitCheckout(branch, ref string) error {
	cmd := exec.Command("git", "checkout", "-b", branch, ref)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitPush(branch string) error {
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

func cleanWorkingDirectory() error {
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
