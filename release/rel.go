package release

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	const sprint = "2510"
	fmt.Printf("Creating %s Branch\n", sprint)

	if len(os.Args) < 2 {
		log.Fatal("new_release_branch_name argument is required")
	}

	lastTag := "2501.0.0-RC.1"

	if lastTag != "" {
		fmt.Printf("Latest and greatest tag is: %s\n", lastTag)
		remoteBranch, err := gitLsRemote(sprint)
		if err != nil {
			log.Fatalf("Error checking remote main: %v", err)
		}
		if remoteBranch != "" {
			fmt.Printf("Remote main %s already exists, hence skipping main creation!!!\n", sprint)
		} else {
			fmt.Printf("Creating the main - %s\n", sprint)
			if err := gitCheckout(sprint, "tags/"+lastTag); err != nil {
				log.Fatalf("Error creating main: %v", err)
			}
			if err := gitPush(sprint); err != nil {
				log.Fatalf("Error pushing main: %v", err)
			}
		}
	} else {
		fmt.Printf("Since there is no tag created for the sprint - %s, cutting the release main based out of master main\n", sprint)
		fmt.Printf("Creating the main - %s\n", sprint)
		if err := gitCheckout(sprint, "master"); err != nil {
			log.Fatalf("Error creating main: %v", err)
		}
		if err := gitPush(sprint); err != nil {
			log.Fatalf("Error pushing main: %v", err)
		}
	}
}

func gitLsRemote(branch string) (string, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", "origin", branch)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
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
