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
		lastTag := LatestTagFetch("RC", setupConfig.current_week_release)
		fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

		incrementTag(lastTag)

		//logic here

	} else {
		//if it does not exist, then creating a new branch for that weekly release
		latestSprintTag := "release." + sprint

		fmt.Printf("Since there is no tag created for the sprint - %s, cutting the release main based out of master main\n", sprint)
		fmt.Printf("Creating the main - %s\n", sprint)
		if err := gitCheckout(latestSprintTag, "origin/main"); err != nil {
			log.Fatalf("Error creating main: %v", err)
		}
		if err := gitPush(latestSprintTag); err != nil {
			log.Fatalf("Error pushing main: %v", err)
		}
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
