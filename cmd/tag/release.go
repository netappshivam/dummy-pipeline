package tag

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Command to handle release creation logic",
	RunE: func(cmd *cobra.Command, args []string) error {
		Release_creation()
		return nil
	},
}

func Release_creation() {
	sprint := SetupConfigobject.CurrentWeekRelease
	fmt.Printf("Checking if branch exists for release-->", sprint)

	// finding if a branch exists for a weekly release
	FetchTagsPrune()
	last_branch, err := FetchReleaseBranch(sprint)
	if err != nil {
		log.Fatalf("Error fetching release branch: %v", err)
	}

	//if it does exists, then finding the latest branch tag and incrementing it
	if last_branch != "" {
		err := CleanWorkingDirectory()
		if err != nil {
			log.Fatalf("Error cleaning working directory: %v", err)
		}

		err2 := GitOnlyCheckout("release." + sprint)
		if err2 != nil {
			log.Fatalf("Error checking out release branch: %v", err2)
		}

		lastTag := LatestTagFetch("RC")
		fmt.Printf("Latest and greatest tag is: %s\n", lastTag)

		newTag, _ := IncrementTag(lastTag)
		fmt.Printf("New tag is: %s\n", newTag)
		fmt.Printf("::set-output name=new_tag::%s\n", newTag)

	} else {
		//if it does not exist, then creating a new branch for that weekly release
		err := CleanWorkingDirectory()
		if err != nil {
			log.Fatalf("Error cleaning working directory: %v", err)
		}

		latestSprintTag := "release." + sprint

		fmt.Printf("Since there is no tag created for the sprint - %s, cutting the release main based out of master main\n", sprint)
		fmt.Printf("Creating the main - %s\n", sprint)
		if err := GitCheckout(latestSprintTag, "origin/main"); err != nil {
			log.Fatalf("Error creating main: %v", err)
		}
		if err := GitPush(latestSprintTag); err != nil {
			log.Fatalf("Error pushing main: %v", err)
		}
		newTag := sprint + ".0.0-RC.1"
		fmt.Printf("::set-output name=new_tag::%s\n", newTag)
	}
}

func init() {
	// Register the releaseCmd with the parent command
	TagCmd.AddCommand(releaseCmd)
}
