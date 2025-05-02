package tag

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Command to handle release creation logic",
	RunE: func(cmd *cobra.Command, args []string) error {
		ReleaseType()
		return nil
	},
}

func ReleaseType() {
	ReleaseFunc()
	ReleaseGithub()
}

func Release_creation() {
	sprint := strings.Split(SetupConfigobject.FinalRelease, ".")[0]
	fmt.Printf("Checking if branch exists for release --> %s\n", sprint)

	// finding if a branch exists for a weekly release
	err_fetch := FetchTagsPrune()
	if err_fetch != nil {
		log.Printf("Failed to fetch tags: %v", err_fetch)
		return
	}

	BranchCheck, err1 := FetchReleaseBranch(sprint)
	if err1 != nil {
		log.Fatalf("Error fetching release branch: %v", err1)
	}

	//if it does not exist, then creating a new branch for that weekly release
	if BranchCheck == "" {
		err := CleanWorkingDirectory()
		if err != nil {
			log.Fatalf("Error cleaning working directory: %v", err)
		}

		latestSprintTag := "release." + sprint

		fmt.Printf("Creating the main - %s\n", sprint)
		if err := GitCheckout(latestSprintTag, SetupConfigobject.BaseRelease); err != nil {
			log.Fatalf("Error creating main: %v", err)
		}
		if err := GitPush(latestSprintTag); err != nil {
			log.Fatalf("Error pushing main: %v", err)
		}

		errCreateGit := CreateGitTag(SetupConfigobject.FinalRelease, "")
		if errCreateGit != nil {
			log.Fatalf("Error creating tag: %v", errCreateGit)
		}

		errGitPush := GitPush(SetupConfigobject.FinalRelease)
		if errGitPush != nil {
			log.Fatalf("Error pushing new tag: %v", errGitPush)
		}

	} else {
		log.Printf("Branch exists")
	}
}

func ReleaseFunc() {
	CurrentTag := FetchDevTag()

	sprint := strings.Split(CurrentTag, ".")[0]

	fmt.Printf("Checking if branch exists for release --> %s\n", sprint)

	// finding if a branch exists for a weekly release
	err_fetch := FetchTagsPrune()
	if err_fetch != nil {
		log.Printf("Failed to fetch tags: %v", err_fetch)
		return
	}

	BranchCheck, err1 := FetchReleaseBranch(sprint)
	if err1 != nil {
		log.Fatalf("Error fetching release branch: %v", err1)
	}

	//if it does not exist, then creating a new branch for that weekly release
	if BranchCheck == "" {
		err := CleanWorkingDirectory()
		if err != nil {
			log.Fatalf("Error cleaning working directory: %v", err)
		}

		devTag := sprint + ".0.0-DEV.1"

		errGitTag := CreateGitTag(devTag, CurrentTag)
		if errGitTag != nil {
			log.Fatalf("Error creating git tag: %v", errGitTag)
		}
		errGitDevPush := GitPush(devTag)
		if errGitDevPush != nil {
			log.Fatalf("Error pushing git tag: %v", errGitDevPush)
		}

		latestSprintBranch := "release/" + sprint

		fmt.Printf("Creating the main - %s\n", sprint)
		if err := GitCheckout(latestSprintBranch, CurrentTag); err != nil {
			log.Fatalf("Error creating main: %v", err)
		}
		if err := GitPush(latestSprintBranch); err != nil {
			log.Fatalf("Error pushing main: %v", err)
		}

		rcTag := sprint + ".0.0-RC.1"

		errCreateGit := CreateGitTag(rcTag, "")
		if errCreateGit != nil {
			log.Fatalf("Error creating tag: %v", errCreateGit)
		}

		errGitPush := GitPush(rcTag)
		if errGitPush != nil {
			log.Fatalf("Error pushing new tag: %v", errGitPush)
		}

		if errWrite := os.WriteFile(os.Getenv("GITHUB_OUTPUT"), []byte(fmt.Sprintf("SAMPLE_RC_TAG=%s\n", SetupConfigobject.FinalRelease)), 0644); errWrite != nil {
			log.Fatalf("Error writing to stdout: %v", errWrite)
		}

	} else {
		log.Printf("Branch exists")
	}
}

func init() {
	// Register the releaseCmd with the parent command
	TagCmd.AddCommand(releaseCmd)
}
