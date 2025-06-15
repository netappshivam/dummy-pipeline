package tag

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var promotionCmd = &cobra.Command{
	Use:   "promotional",
	Short: "Command to handle release creation logic",
	RunE: func(cmd *cobra.Command, args []string) error {
		SetupConfigobject.PromotionalFunc()
		return nil
	},
}

func (o *SetupConfig) PromotionalFunc() {
	if o.OperationType == "Final_Tag" {
		errFetch := FetchTagsPrune()
		if errFetch != nil {
			log.Printf("Failed to fetch tags: %v", errFetch)
			return
		}
		log.Printf("base: %s", o.BaseRelease)
		log.Printf("final: %s", o.FinalRelease)
		log.Printf("Operation Type: %s", o.OperationType)

		a := strings.Contains(o.BaseRelease, "-RC.")
		log.Printf("Value of a: %v", a)

		b := strings.Contains(o.BaseRelease, o.FinalRelease)
		log.Printf("Value of b: %v", b)

		c := FetchTag(o.FinalRelease)
		cname := false
		if c == "" {
			log.Printf("Tag %s does not exist", o.FinalRelease)
			cname = true
		} else {
			log.Printf("Tag %s exists", o.FinalRelease)
		}

		d := FetchTag(o.BaseRelease)
		dname := false
		if d == "" {
			log.Printf("Tag %s does not exist", o.BaseRelease)
		} else {
			log.Printf("Tag %s exists", o.BaseRelease)
			dname = true
		}

		if dname && strings.Contains(o.BaseRelease, "-RC.") && o.FinalRelease == o.BaseRelease[:10] && cname {
			log.Println("Base release is a valid RC tag, proceeding with promotion.")
			PromotionalCreation()
		} else {
			log.Fatalf("Base release %s is not a valid RC tag or Final Release has a naming error %s", o.BaseRelease, o.FinalRelease)
			return
		}
	} else if o.OperationType == "HF_Release" {
		if FetchTagToCheckIfItExists(o.BaseRelease) && CheckForHFfinalName(o) && !FetchTagToCheckIfItExists(o.FinalRelease) {
			log.Println("Base release is a valid HF tag, proceeding with promotion.")
			PromotionalCreation()
		} else {
			log.Fatalf("Base release %s is not a valid HF tag or Final Release has a naming error %s", o.BaseRelease, o.FinalRelease)
			return
		}
	}
}

func PromotionalCreation() {
	errFetch := FetchTagsPrune()
	if errFetch != nil {
		log.Printf("Failed to fetch tags: %v", errFetch)
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

	if errWrite := os.WriteFile(os.Getenv("GITHUB_OUTPUT"), []byte(fmt.Sprintf("FINAL_TAG=%s\n", SetupConfigobject.FinalRelease)), 0644); errWrite != nil {
		log.Fatalf("Error writing to stdout: %v", errWrite)
	}
}

func init() {
	TagCmd.AddCommand(promotionCmd)
}
