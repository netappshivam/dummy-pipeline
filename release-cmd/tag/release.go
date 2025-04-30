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
	if strings.Contains(SetupConfigobject.BaseRelease, "-DEV.") {
		Release_creation()
	} else if strings.Contains(SetupConfigobject.BaseRelease, "-RC.") {
		PromotionalCreation()
	}
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

func PromotionalCreation() {
	err_fetch := FetchTagsPrune()
	if err_fetch != nil {
		log.Printf("Failed to fetch tags: %v", err_fetch)
		return
	}

	err := CleanWorkingDirectory()
	if err != nil {
		log.Fatalf("Error cleaning working directory: %v", err)
	}

	errGitTag := CreateGitTag(SetupConfigobject.FinalRelease, SetupConfigobject.BaseRelease)
	if errGitTag != nil {
		log.Fatalf("Error creating git tag: %v", errGitTag)
	}

	errGitPush := GitPush(SetupConfigobject.FinalRelease)
	if errGitPush != nil {
		log.Fatalf("Error pushing git tag: %v", errGitPush)
	}

	errSetEnv := os.Setenv("FINAL_RELEASE_TAG", SetupConfigobject.FinalRelease)
	if errSetEnv != nil {
		log.Fatalf("Error setting environment variable: %v", errSetEnv)
	}

}

func init() {
	// Register the releaseCmd with the parent command
	TagCmd.AddCommand(releaseCmd)
}
